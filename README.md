# league-api
GRPC API for league stats

- Full GRPC server for a simple API that connects to the Riot API
- Uses the grpc-gateway to create a rest api
- Adds a simple rest api just for fun and learning.

# Notes
sudo cp -a Downloads/include /usr/local/include

#Generate GRPC stub

- Standard Protobuf

protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. leagueservice/league.proto

- GOGO Protobuf

protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gofast_out=plugins=grpc:. leagueservice/league.proto

- Generates Reverse Proxy

protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. leagueservice/league.proto 

# TODO
- Add Docker file maybe K8s
- Redis for caching profile calls
- Add channels and go routines for the api calls
- Try to get rid of intermediate structs
- Restructure the rest api and client app to best practices
