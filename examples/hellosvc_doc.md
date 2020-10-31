# A Practice for generate microservice [hello]

It contains some features:

- gRPC transport
- Customize pb dir
- import proto file from other dir

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
$ mkdir hello/pb
```

Then execute the kit command:
```bash
$ pwd
kit/examples

$ kit g s hello -t grpc -p hello/pb -i hello/pb
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

## 3. Restructure pb directory
In order to be as consistent as possible with the actual environment, now we restruture pb dir

Create subfolder `proto` in `/hello/pb/`:
```bash
$ mkdir hello/pb/proto 
```
Move all proto files and scripts into `proto` dir, create a subfolder `common` to put sub protobuffer dir,
this dir used to store common proto file which would be imported from proto files within `proto` dir.
 
now the structure like following:
```bash
examples/hello/pb/proto/
│  compile.sh
│  hello.proto
│
└─common
```
Create common proto file `common.proto` in `hello/pb/proto/common/`
```bash
$ vi hello/pb/proto/common/common.proto
syntax = "proto3";

package pb;
option go_package = "hello/pb/gen-go/common;common";

message BaseReq {
	string token=1;
}
```

## 3. Separate proto files from `*.pb.go` files

You may have noticed `option go_package` at `common.proto`, yes, we will put `.pb.go` files into 
`pb/gen-go` dir, but no need to create it manually.

Now we modify the `option go_package` at `hello.proto`:
```proto
// old
option go_package = "hello/pb;pb";

// new
option go_package = "hello/pb/gen-go/pb;pb";
``` 


Then import `common.BaseReq` at `hello.proto`:
```protobuf
// hello.proto
syntax = "proto3";
package pb;

import "common/common.proto";  <---------- use relative path
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
protoc *.proto --go_out=plugins=grpc:../../../
protoc common/*.proto --go_out=plugins=grpc:../../../   <----------- new line
```

If you want to execute `compile.bat` script with a subfolder stored more proto
files on Cmd(windows), see also `hello/pb/proto/compile.bat`.

## Regenerate pb files
```bash
$ pwd
kit/examples

# note: -i option has been changed to hello/pb/gen-go/pb
$ kit g s hello -t grpc -p hello/pb/proto -i hello/pb/gen-go/pb
time="2020-09-19T08:32:29+08:00" level=info msg="exec>[sh -c /.../kit/examples/hello/pb/proto/compile.sh]"
time="2020-09-19T08:27:50+08:00" level=warning msg="==============================================================="
time="2020-09-19T08:27:50+08:00" level=warning msg="The GRPC implementation is not finished you need to update your"
time="2020-09-19T08:27:50+08:00" level=warning msg=" service proto buffer and run the compile script."
time="2020-09-19T08:27:50+08:00" level=warning msg=---------------------------------------------------------------
time="2020-09-19T08:27:50+08:00" level=warning msg="You also need to implement the Encoders and Decoders!"
time="2020-09-19T08:27:50+08:00" level=warning msg="==============================================================="
```

This will generate `.pb.go` files in the location specified by `option go_package` in proto file,
for here, it's `hello/pb/gen-go/`. If the location is not correspond to `option go_package` in proto file,
you need to check compile script, just refer to the method mentioned above. 
 
You can check whether the `.go` file import path in following dir is correct:
- `hello/pkg/grpc/`
- `hello/cmd/service/`

if not, report to us via Issue.

## 4. Run it
```bash
$ cd hello/cmd/
$ go run .
ts=2020-09-18T07:57:07.021432Z caller=service.go:84 tracer=none
ts=2020-09-18T07:57:07.022408Z caller=service.go:106 transport=gRPC addr=:8082
ts=2020-09-18T07:57:07.022408Z caller=service.go:134 transport=debug/HTTP addr=:8080
```

Anything abnormally about kit, feel free to report to us, enjoy kit~