#!/bin/bash
set -x

# register
curl -X POST http://localhost:8800/user/register  -H "Content-Type: application/json" -d '{"first_name":"Ivan10","second_name":"Frolov10","birthdate":"2002-02-11","biography":"Hokkey","city":"Moskva","password":"passI"}'

# login
curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"fb0e9288-e4d2-9561-166a-eb34120bc3c3","password":"passI"}'


#get user
curl -i http://localhost:8800/user/get/fb0e9288-e4d2-9561-166a-eb34120bc3c3



