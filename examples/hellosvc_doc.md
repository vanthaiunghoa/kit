# Practice generating microservice [hello]

It contained some features:

- With gRPC transport
- Customize pb dir
- import proto file

## 1. Create a service
```bash
$kit n s hello
```

Then add a method for service, edit `hello/pkg/service/service.go`
```go
// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	Foo(ctx context.Context, s string) (rs string, err error)
}
```

## 2. Generate transport and endpoint layer

Create your pb dir firstly:
```bash
$mkdir hello/pb
```

Then execute the kit command:
```bash
$pwd
kit/examples

$kit g s hello -t grpc -p hello/pb -i hello/pb
time="2020-09-18T10:23:36+08:00" level=warning msg="==============================================================="
time="2020-09-18T10:23:36+08:00" level=warning msg="The GRPC implementation is not finished you need to update your"
time="2020-09-18T10:23:36+08:00" level=warning msg=" service proto buffer and run the compile script."
time="2020-09-18T10:23:36+08:00" level=warning msg=---------------------------------------------------------------
time="2020-09-18T10:23:36+08:00" level=warning msg="You also need to implement the Encoders and Decoders!"
time="2020-09-18T10:23:36+08:00" level=warning msg="==============================================================="
```

With here, `-p` set proto file path, `-i` set import path of pb file.
If you haven't change default proto file path via `-p`, there is no need to set `-i` too.

It can be checked in `pkg/grpc/handler.go`.


There is one more thing need to note, after you modified proto option `go_package` with a more detailed path, 
like `go_package="hello/pb;pb"`, since it is full go pkg path startswith project root dir, 
so you also need modify `compile.sh` like following:

```bash
// old
protoc hello.proto --go_out=plugins=grpc:.

// new
protoc hello.proto --go_out=plugins=grpc:../../
```

Otherwise, The location of the generated pb file will be out of order, this is just a path issue.

## 3. Add a new proto file

We add a new proto file now:
```bash
$vi hello/pb/common.proto
syntax = "proto3";
package pb;

option go_package = "hello/pb;pb";
message BaseReq {
	string token=1;
}
``` 
Then import `BaseReq` from `hello.proto`:
```protobuf
// hello.proto
syntax = "proto3";
package pb;

import "common.proto"; <----------
option go_package = "hello/pb;pb";

//The Hello service definition.
service Hello {
	  rpc Foo (FooRequest) returns (FooReply);
}
message FooRequest {
	  BaseReq base_req = 1;   <----------
}
message FooReply {}
```

Because we added a new proto file, we need modify `compile.sh`:
```bashtext
// old
protoc hello.proto --go_out=plugins=grpc:.

// new, notice *.proto is not supported on cmd env(windows).
protoc *.proto --go_out=plugins=grpc:../../
```

If you want to execute `compile.bat` script with a subfolder stored more proto
files on Cmd(windows), see also example `hello/pb/compile.bat`


## 4. Run it
```bash
$cd hello/cmd/
$go run .
ts=2020-09-18T07:57:07.021432Z caller=service.go:84 tracer=none
ts=2020-09-18T07:57:07.022408Z caller=service.go:106 transport=gRPC addr=:8082
ts=2020-09-18T07:57:07.022408Z caller=service.go:134 transport=debug/HTTP addr=:8080
ts=2020-09-18T07:57:08.0354949Z caller=service.go:93 exit="received signal interrupt"
```