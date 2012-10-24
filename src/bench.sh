#!/bin/sh

# ab uses HTTP/1.0 causing an extra header for Go (Connection: close). Java
# violates HTTP by sending a 1.1 response to a 1.0 request!
#ab -A apikey_value -c 1 -n 1 http://localhost:8080/authenticate

# siege sends an Accept-Encoding: gzip header, which causes Java to gzip the
# response, but Go does not.
#siege -b -r 1 -c 1 -H 'Authorization: basic apikey_value' http://localhost:8080/authenticate

# httperf causes identical behavior at the HTTP level
httperf --hog --add-header 'Authorization: basic apikey_value\n' \
  --server localhost --port 8080 --uri /authenticate \
  --num-calls 1000 --num-conns 50

curl -i -H "Authorization: uniqush rsaencryptedclientkey" http://localhost:8080/authenticate