#!/bin/bash
set -x

# register
curl -X POST http://localhost:8800/user/register  -H "Content-Type: application/json" -d '{"first_name":"Ivan10","second_name":"Frolov10","birthdate":"2002-02-11","biography":"Hokkey","city":"Moskva","password":"passI"}'

# login
curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"28832b26-ef3e-f373-f318-fbce83b9c90f","password":"password"}'

curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"9d3a433d-d242-b1e2-d9f2-dfcd9a248638","password":"password"}'


#get user
curl -i http://localhost:8800/user/get/28832b26-ef3e-f373-f318-fbce83b9c90f

#get feed posts
curl http://localhost:8800/post/feed  -H "Authorization: Bearer B4WL7S3DM2RSD2Q4KPOKO5QYAQ"

#post create
curl http://localhost:8800/post/create  -H "Authorization: Bearer 5IPEH4YNR3EPJTL5QRV42MYD5E" -d '{"text":"text6"}'

curl -i -X PUT http://localhost:8800/friend/set/654ae314-e325-5909-4e1e-1f694b00b699  -H "Authorization: Bearer 5IPEH4YNR3EPJTL5QRV42MYD5E"




curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"9e8dcc01-abb2-2502-d900-f0bd0074754d","password":"password"}'
curl http://localhost:8800/post/feed  -H "Authorization: Bearer ETEYF6C3ERPVPBCNGC6X6AP2CY"


curl http://localhost:8800/post/feed?limit=100  -H "Authorization: Bearer ETEYF6C3ERPVPBCNGC6X6AP2CY"

curl http://localhost:8800/postsupdate

curl http://localhost:8800/post/feed?limit=100  -H "Authorization: Bearer ETEYF6C3ERPVPBCNGC6X6AP2CY"
