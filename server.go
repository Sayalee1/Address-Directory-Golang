package main

import (
  "fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "github.com/julienschmidt/httprouter"
  "net/http"
  "encoding/json"
  "strings"
)

type LocReq struct {
    Name    string `json:"name"`
    Address string `json:"address"`
    City    string `json:"city"`
    State   string `json:"state"`
    Zip     string `json:"zip"`
}

type LocResp struct {
    ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Name  string `json:"name"`
    Address    string `json:"address"`
    City       string `json:"city"`
    State string `json:"state"`
    Zip   string `json:"zip"`
    Coordinate struct {
        Lat float64 `json:"lat"`
        Lng float64 `json:"lng"`
    } `json:"coordinate"`
}
type GeoLocation struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PlaceID string   `json:"place_id"`
        Types   []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}

type GoogleLocResp struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PlaceID string   `json:"place_id"`
        Types   []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}
var c *mgo.Collection
var locResp LocResp

const(
    timeout = time.Duration(time.Second*100)
)
func dbConnection() {
  uri := "mongodb://admin:admin@ds045970.mongolab.com:45970/location"
    session, err := mgo.Dial(uri)
    if err != nil {
      panic(err)
    } else {
      session.SetSafe(&mgo.Safe{})
      c = session.DB("location").C("location")
  }
}

func main() {
    mux := httprouter.New()
    mux.GET("/locations/:id", getLoc)
    mux.POST("/locations", addLoc)
    mux.PUT("/locations/:id", updateLoc)
    mux.DELETE("/locations/:id", delLoc)
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }
  dbConnection()
    server.ListenAndServe()
}
func getLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
  
  id := bson.ObjectIdHex(p.ByName("id"))
  err := c.FindId(id).One(&locResp)
    if err != nil {
      panic(err)
    }   
  rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
    rw.WriteHeader(200)
    json.NewEncoder(rw).Encode(locResp)
}
func getGoogleLocation(address string) (gLoc GeoLocation) {
  
  client := http.Client{Timeout: timeout}
  googleurl := fmt.Sprintf("http://maps.google.com/maps/api/geocode/json?address=%s",address)
    res, err := client.Get(googleurl)
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()

    decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&gLoc)
    if(err!=nil)    {
        panic(err)
    } 
  return gLoc
}
func addLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    var lReq LocReq
    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&lReq)
    if(err!=nil)    {
        panic(err)
    }
  address := lReq.Address+" "+lReq.City+" "+lReq.State+" "+lReq.Zip
  address = strings.Replace(address," ","%20",-1)

  locationDetails := getGoogleLocation(address)

    locResp.ID = bson.NewObjectId()
  locResp.Address= lReq.Address
  locResp.City=lReq.City
  locResp.Name=lReq.Name
  locResp.State=lReq.State
  locResp.Zip=lReq.Zip
  locResp.Coordinate.Lat=locationDetails.Results[0].Geometry.Location.Lat
  locResp.Coordinate.Lng=locationDetails.Results[0].Geometry.Location.Lng

  err = c.Insert(locResp)
  if err != nil {
      panic(err)
    }

  err = c.FindId(locResp.ID).One(&locResp)
    if err != nil {
      panic(err)
    }   
  rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
    rw.WriteHeader(201)
    json.NewEncoder(rw).Encode(locResp)
}

func updateLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

  var tempLocResp LocResp
  var locResp LocResp
  id := bson.ObjectIdHex(p.ByName("id"))
  err := c.FindId(id).One(&locResp)
    if err != nil {
      panic(err)
    } 
  tempLocResp.Name = locResp.Name
  tempLocResp.Address = locResp.Address
  tempLocResp.City = locResp.City
  tempLocResp.State = locResp.State
  tempLocResp.Zip = locResp.Zip
    decoder := json.NewDecoder(req.Body)
    err = decoder.Decode(&tempLocResp)
  
    if(err!=nil)    {
        panic(err)
    }

  address := tempLocResp.Address+" "+tempLocResp.City+" "+tempLocResp.State+" "+tempLocResp.Zip
  address = strings.Replace(address," ","%20",-1)
  locationDetails := getGoogleLocation(address)
  tempLocResp.Coordinate.Lat=locationDetails.Results[0].Geometry.Location.Lat
  tempLocResp.Coordinate.Lng=locationDetails.Results[0].Geometry.Location.Lng
  err = c.UpdateId(id,tempLocResp)
    if err != nil {
      panic(err)
    } 

  err = c.FindId(id).One(&locResp)
    if err != nil {
      panic(err)
    }
  rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
    rw.WriteHeader(201)
    json.NewEncoder(rw).Encode(locResp)
}

func delLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

  id := bson.ObjectIdHex(p.ByName("id"))
  err := c.RemoveId(id)
    if err != nil {
      panic(err)
    }
  rw.WriteHeader(200)
}