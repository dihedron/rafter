.PHONY: binary
binary:
	@go build

.PHONY: clean
clean:
	@go clean

.PHONY: reset
reset:
	@go clean -x -cache

.PHONY: proto
proto: proto/service.proto
	@rm -f proto/service.pb.go proto/service_grpc.pb.go
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto

.PHONY: run
run: binary
	@if [ -d "tests/raft/store" ]; then \
		echo "running existing cluster..."; \
		goreman -f tests/raft/run.procfile start; \
	else \
		echo "bootstrapping new cluster..."; \
		goreman -f tests/raft/bootstrap.procfile start; \
	fi

.PHONY: bootstrap
bootstrap:
	@rm -rf tests/raft/store
	@make run