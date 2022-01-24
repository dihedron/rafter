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
proto: application/proto/service.proto
	@rm -f application/proto/service.pb.go application/proto/service_grpc.pb.go
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative application/proto/service.proto

.PHONY: run3
run3: binary
	@if [ -d "tests/raft/store" ]; then \
		echo "running existing cluster..."; \
		goreman -f tests/raft/run3.procfile start; \
	else \
		echo "bootstrapping new cluster..."; \
		goreman -f tests/raft/bootstrap3.procfile start; \
	fi

.PHONY: bootstrap3
bootstrap3:
	@rm -rf tests/raft/store
	@make run3

.PHONY: benchmark3
benchmark3:
	@./rafter data benchmark --length=64 --iterations=1000 --concurrency=10 --peer=@tests/raft/node1.json --peer=@tests/raft/node2.json --peer=@tests/raft/node3.json  

.PHONY: run5
run5: binary
	@if [ -d "tests/raft/store" ]; then \
		echo "running existing cluster..."; \
		goreman -f tests/raft/run5.procfile start; \
	else \
		echo "bootstrapping new cluster..."; \
		goreman -f tests/raft/bootstrap5.procfile start; \
	fi

.PHONY: bootstrap5
bootstrap5:
	@rm -rf tests/raft/store
	@make run5

.PHONY: benchmark5
benchmark5:
	@./rafter data benchmark --length=64 --iterations=1000 --concurrency=10 --peer=@tests/raft/node1.json --peer=@tests/raft/node2.json --peer=@tests/raft/node3.json --peer=@tests/raft/node4.json --peer=@tests/raft/node5.json 
