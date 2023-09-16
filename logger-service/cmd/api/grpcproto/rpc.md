go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
create proto file
install protoc:  sudo apt install protobuf-compiler
Run the following command from the directory that has the protobuf file:
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative filename
make sure proto-gen-go is present in class path

-> The command will generate two files filename.pb.go and filename_grpc.pb.go
-> Add go get google.golang.org/grpc
   go get google.golang.org/protobuf