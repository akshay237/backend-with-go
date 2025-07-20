## GRPC
* It is a remote procedure call framework where 
- the client can execute the remote procedure call on server.
- the api and data structures code is automatically generated.
- Supports multiple programming language.

### How it works:

* Define API and Data Structures
    - The RPC and it's request/response structure are defined using protobuf.
* Generate gRPC stubs
    - Generates codes for the server and client in the language of choice.
* Implement the server
    - Implement the RPC handler on the server side.
* Use the client
    - Use the generated client stubs to call the RPC on server

### Why GRPC?

* High Performance
* Strong API Contract
* Automatic Code Generation

### Types of GRPC

* Unary gRPC
    - Same like normal HTTP request.
* Client Streaming gRPC
    - where client sends stream multiple request and server responds once for all of them
* Server Streaming gRPC
    - Where client sends one request and server responds in streaming to the client
* Bidirectional Streaming gRPC
    - Here client sends streams mutliplr request and server responds with streaming

### gRPC Gateway:

* Serve both GRPC and HTTP Requests at once
* A plugin of protocol buffer compiler
* Write code once serves both
* Translate HTTP Json calls to gRPC
    - in-process transalation for unary (no extra network hop)
    - seprate proxy: via network call serves both unary and streaming 


