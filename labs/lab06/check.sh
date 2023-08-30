#!/bin/bash

set -x

#friends
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"34e77d29-52dc-3187-e572-ef093cefed7f","password":"password"}'
#{"token":"4BACPA66EHIY5LJTCREDLZR2P4"}
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"460dc67d-3b8c-4b9a-669c-8f28dcb77d84","password":"password"}'
#{"token":"635YFENQLRA7QRH5OAJWLBK5UQ"}



#user
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"e2054f50-f7b3-c48f-d710-d6186fbc4caf","password":"password"}'
#{"token":"33PRO6E2O3TC4APOHIXONNQSYI"}


curl -i http://localhost:8800/post/create  -H "Authorization: Bearer 33PRO6E2O3TC4APOHIXONNQSYI" -d '{"text":"post1 from user e2054f50"}'
curl -i http://localhost:8800/post/create  -H "Authorization: Bearer 33PRO6E2O3TC4APOHIXONNQSYI" -d '{"text":"post2 from user e2054f50"}'
curl -i http://localhost:8800/post/create  -H "Authorization: Bearer 33PRO6E2O3TC4APOHIXONNQSYI" -d '{"text":"post3 from user e2054f50"}'
