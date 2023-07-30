#!/bin/bash
set -x

# register
curl -X POST http://localhost:8800/user/register  -H "Content-Type: application/json" -d '{"first_name":"Ivan10","second_name":"Frolov10","birthdate":"2002-02-11","biography":"Hokkey","city":"Moskva","password":"passI"}'

# login
curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"37b48b26-03b1-9d29-57ad-b306b41edbda","password":"passI"}'


#get user
curl -i http://localhost:8800/user/get/37b48b26-03b1-9d29-57ad-b306b41edbda


