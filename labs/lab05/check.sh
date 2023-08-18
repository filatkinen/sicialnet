#!/bin/bash

set -x

#shard 1 users:
#628ad2a9-1676-5aba-f27c-dd38363128b8,Екатерина,Артемова,,,Октябрьский,2001-01-02 00:00:00.000000,1
#902c7d85-ace2-319d-6852-4224b2866067,Яна,Артемова,,,Тимашевск,2001-01-14 00:00:00.000000,1
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"628ad2a9-1676-5aba-f27c-dd38363128b8","password":"password"}'
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"902c7d85-ace2-319d-6852-4224b2866067","password":"password"}'
#{"token":"6366JCTUED2OR6FGXRRMWGRXWQ"}
#{"token":"KBJHYY7HLYUB6BP7YQOJC4VLIE"}

#shard 2 users:
#aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b,Макар,Беляев,,,Майкоп,2001-01-04 00:00:00.000000,2
#4090b741-73e9-17b2-909c-d89913f33fc3,Арина,Беляева,,,Нальчик,2001-01-02 00:00:00.000000,2
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b","password":"password"}'
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b","password":"password"}'
#{"token":"SS7KA4B6VL4AU2D5DDMV36VLWQ"}
#{"token":"UR323CUFFHRZZ6MIBKUOROPLII"}

#
#shard 3 users:
#25266b4f-0b79-e2a5-8cc3-c2f772fc43ec,Егор,Беликов,,,Архангельск,2001-01-05 00:00:00.000000,3
#45c497c7-d441-3c6c-eeb9-bf34c819ad43,Василиса,Белоусова,,,Мичуринск,2001-01-03 00:00:00.000000,3
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"25266b4f-0b79-e2a5-8cc3-c2f772fc43ec","password":"password"}'
#curl -X POST http://localhost:8800/login -H "Content-Type: application/json" -d '{"id":"45c497c7-d441-3c6c-eeb9-bf34c819ad43","password":"password"}'
#{"token":"IAAEYSFMW2D4RGWBIEKROKTWWY"}
#{"token":"VUZQRGVDAATLOQPKEIE5A5Y6FU"}



####Shard 2
#send user shard 1 to user shard2. We are to see message in shard2
curl -i -X POST http://localhost:8800/dialog/aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b/send  -H "Authorization: Bearer 6366JCTUED2OR6FGXRRMWGRXWQ" -d '{"text":"message 0101"}'
curl -i http://localhost:8800/dialog/aaaf8f6f-f7c7-9e42-d01d-c10a35391a7b/list  -H "Authorization: Bearer 6366JCTUED2OR6FGXRRMWGRXWQ"

#send user shard 2 to user shard1. We are to see message in shard2
curl -i -X POST http://localhost:8800/dialog/628ad2a9-1676-5aba-f27c-dd38363128b8/send  -H "Authorization: Bearer SS7KA4B6VL4AU2D5DDMV36VLWQ" -d '{"text":"message 0101"}'
curl -i http://localhost:8800/dialog/628ad2a9-1676-5aba-f27c-dd38363128b8/list  -H "Authorization: Bearer SS7KA4B6VL4AU2D5DDMV36VLWQ"


####Shard 3
#send user shard 3 to user shard3. We are to see message in shard3
curl -i -X POST http://localhost:8800/dialog/45c497c7-d441-3c6c-eeb9-bf34c819ad43/send  -H "Authorization: Bearer IAAEYSFMW2D4RGWBIEKROKTWWY" -d '{"text":"message 0101"}'
curl -i http://localhost:8800/dialog/45c497c7-d441-3c6c-eeb9-bf34c819ad43/list  -H "Authorization: Bearer IAAEYSFMW2D4RGWBIEKROKTWWY"

#send user shard 3 to user shard3. We are to see message in shard2
curl -i -X POST http://localhost:8800/dialog/25266b4f-0b79-e2a5-8cc3-c2f772fc43ec/send  -H "Authorization: Bearer VUZQRGVDAATLOQPKEIE5A5Y6FU" -d '{"text":"message 0102"}'
curl -i http://localhost:8800/dialog/25266b4f-0b79-e2a5-8cc3-c2f772fc43ec/list  -H "Authorization: Bearer VUZQRGVDAATLOQPKEIE5A5Y6FU"

