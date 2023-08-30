#!/bin/bash

set -x


docker exec -it socialnet_wsclient /usr/src/app/socialnet/bin/wsclient -token 635YFENQLRA7QRH5OAJWLBK5UQ -url 'ws://socialnet_app:8800/post/feed/posted'
