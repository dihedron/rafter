log.pb.go: log.proto
#	protoc service.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative log.proto

force:
	rm -f service.pb.go
	make service.pb.go
