This is a test script for a test deployment on AWS ECS with an ALB that's supposed to be http2 both ways.

AWS is running this [grpc service](https://github.com/grpc/grpc-java/tree/master/examples/example-hostname) (docker
image available at: `grpc/java-example-hostname`) at `grpc.002fa7.net`, run the client with:

```json
go generate
go run .
```

In this commit's current state, we error out with:

```
2021/10/08 23:20:56 error in request: rpc error: code = Unavailable desc = connection closed before server preface received
panic: error in request: rpc error: code = Unavailable desc = connection closed before server preface received
```
