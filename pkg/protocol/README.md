To generate protobuf files run;

```shell
go install google.golang.org/protobuf/cmd/pro toc-gen-go@v1.28.1

cd pkg/protocol/protobuf
protoc -I . --go_out=. --go_opt=module=github.com/teslamotors/vehicle-command/pkg/protocol/protobuf *.proto
```
