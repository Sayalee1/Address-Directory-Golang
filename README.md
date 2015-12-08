# cmpe273-assignment2
The location service has the following REST endpoints to store and retrieve locations. All the data persists into MongoDB(mongolab).

Create New Location - POST        /locations
Request:
POST /locations
{
   "name" : "Sayalee Agashe",
   "address" : "123 Main St",
   "city" : "San Francisco",
   "state" : "CA",
   "zip" : "94113"
}
Response
HTTP Response Code: 201
{
   "id" : 12345,
   "name" : "Sayalee Agashe",
   "address" : "123 Main St",
   "city" : "San Francisco",
   "state" : "CA",
   "zip" : "94113",
   "coordinate" : { 
      "lat" : 38.4220352,
     "lng" : -222.0841244
   }
}
Get a Location - GET        /locations/{location_id}
Request
GET /locations/12345
Response
HTTP Response Code: 200
{
   "id" : 12345,
   "name" : "Sayalee Agashe",
   "address" : "123 Main St",
   "city" : "San Francisco",
   "state" : "CA",
   "zip" : "94113",
   "coordinate" : { 
      "lat" : 38.4220352,
     "lng" : -222.0841244
   }
}
Update a Location - PUT /locations/{location_id}
Request:
PUT /locations/12345
{
   "address" : "1600 Amphitheatre Parkway",
   "city" : "Mountain View",
   "state" : "CA",
   "zip" : "94043"
}
Response
HTTP Response Code: 201
{
   "id" : 12345,
   "name" : "Sayalee Agashe",
   "address" : "1600 Amphitheatre Parkway",
   "city" : "Mountain View",
   "state" : "CA",
   "zip" : "94043"
   "coordinate" : { 
      "lat" : 37.4220352,
     "lng" : -122.0841244
   }
}
Delete a Location - DELETE /locations/{location_id}
        Request:
DELETE  /locations/12345
        Response:
HTTP Response Code: 200
