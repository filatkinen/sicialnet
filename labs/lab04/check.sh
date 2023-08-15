#!/bin/bash
set -x


hey -n 10000 -c 1000 -H "Authorization: Bearer ETEYF6C3ERPVPBCNGC6X6AP2CY" -m GET 'http://localhost:8800/post/feed?limit=100'


curl http://localhost:8800/postsupdate

sleep 30

hey -n 10000 -c 1000 -H "Authorization: Bearer ETEYF6C3ERPVPBCNGC6X6AP2CY" -m GET 'http://localhost:8800/post/feed?limit=100'
