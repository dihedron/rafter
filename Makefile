.PHONY: binary
binary:
	@go build

.PHONY: protobuf
protobuf: service.proto
	cd proto
	rm -f service.pb.go service_grpc.pb.go
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative service.proto
	cd..

.PHONY: run
run: binary
	if [ -d "tests/raft/store" ]; then \
		echo "running existing cluster..."; \
		goreman -f tests/raft/run.procfile start; \
	else \
		echo "bootstrapping new cluster..."; \
		goreman -f tests/raft/bootstrap.procfile start; \
	fi

.PHONY: bootstrap
bootstrap:
	rm -rf tests/raft/store
	make run