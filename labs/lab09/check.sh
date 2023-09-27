#!/bin/bash

set -x


#628ad2a9-1676-5aba-f27c-dd38363128b8,Екатерина,Артемова,,,Октябрьский,2001-01-02 00:00:00.000000,1
#902c7d85-ace2-319d-6852-4224b2866067,Яна,Артемова,,,Тимашевск,2001-01-14 00:00:00.000000,1
#curl -X POST http://localhost:8801/login -H "Content-Type: application/json" -d '{"id":"628ad2a9-1676-5aba-f27c-dd38363128b8","password":"password"}'
#{"token":"QFBPIXM7S4VQVPF2TKA5ZBPQII"}
#curl -X POST http://localhost:8801/login -H "Content-Type: application/json" -d '{"id":"902c7d85-ace2-319d-6852-4224b2866067","password":"password"}'
#{"token":"OGJHGZ76N654Z33JYKW65CY3NU"}

curl -i -X POST http://localhost:8801/dialog/902c7d85-ace2-319d-6852-4224b2866067/send  -H "Authorization: Bearer QFBPIXM7S4VQVPF2TKA5ZBPQII" -d '{"text":"message1 From Katya to Anna"}'
curl -i -X POST http://localhost:8801/dialog/902c7d85-ace2-319d-6852-4224b2866067/send  -H "Authorization: Bearer QFBPIXM7S4VQVPF2TKA5ZBPQII" -d '{"text":"message2 From Katya to Anna"}'

curl -i -X POST http://localhost:8801/dialog/628ad2a9-1676-5aba-f27c-dd38363128b8/send  -H "Authorization: Bearer OGJHGZ76N654Z33JYKW65CY3NU" -d '{"text":"message1 From Anna to Katya"}'
curl -i -X POST http://localhost:8801/dialog/628ad2a9-1676-5aba-f27c-dd38363128b8/send  -H "Authorization: Bearer OGJHGZ76N654Z33JYKW65CY3NU" -d '{"text":"message2 From Anna to Katya"}'
