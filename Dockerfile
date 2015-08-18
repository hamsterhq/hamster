FROM golang:1.4.2
RUN go get -u github.com/BurntSushi/toml \
  github.com/adnaan/routes \
  github.com/garyburd/redigo/internal \
  github.com/garyburd/redigo/redis \
  github.com/gorilla/context \
  github.com/gorilla/securecookie \
  github.com/gorilla/sessions \
  github.com/kr/fernet \
  golang.org/x/crypto/blowfish \
  labix.org/v2/mgo \
  labix.org/v2/mgo/bson \
  code.google.com/p/go.crypto/bcrypt
