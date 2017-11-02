
## dependencies following packages
+ Download and install protoc compiler from here: https://github.com/google/protobuf. Add the location of protoc binary 
  file into PATH environment variable so that you can invoke protoc compiler from any location.
+ Install the protoc plugin for your language. For Go, run the go get command to install the protoc plugin for Go:

        go get -u github.com/golang/protobuf/protoc-gen-go
        go get -u github.com/golang/protobuf/proto
    
In order to working with gRPC in Go, you must install Go implementation of gRPC:

    go get google.golang.org/grpc
    
Finally ,build ouput go file with following command:

    protoc -I customer/ customer/customer.proto --go_out=plugins=grpc:customer