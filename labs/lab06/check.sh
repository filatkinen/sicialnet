#!/bin/bash

set -x


curl -i http://localhost:8800/post/create  -H "Authorization: Bearer 33PRO6E2O3TC4APOHIXONNQSYI" -d '{"text":"post1 from user e2054f50"}'
curl -i http://localhost:8800/post/create  -H "Authorization: Bearer 33PRO6E2O3TC4APOHIXONNQSYI" -d '{"text":"post2 from user e2054f50"}'
curl -i http://localhost:8800/post/create  -H "Authorization: Bearer 33PRO6E2O3TC4APOHIXONNQSYI" -d '{"text":"post3 from user e2054f50"}'
