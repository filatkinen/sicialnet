#!/bin/bash

set -x


docker exec -it socialnet_wsclient /usr/src/app/socialnet/bin/wsclient -token 4BACPA66EHIY5LJTCREDLZR2P4 -url 'ws://socialnet_app:8800/post/feed/posted'

