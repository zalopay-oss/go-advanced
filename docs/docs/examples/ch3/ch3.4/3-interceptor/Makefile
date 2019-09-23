
gen:
	protoc -I . --go_out=plugins=grpc:. helloworld.proto

clean:
	-rm *.pb.go *interceptor
