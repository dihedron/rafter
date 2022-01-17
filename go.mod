module github.com/dihedron/rafter

go 1.17

require (
	github.com/Jille/raft-grpc-leader-rpc v1.1.0
	github.com/Jille/raft-grpc-transport v1.2.0
	github.com/Jille/raftadmin v1.2.0
	github.com/dihedron/grpc-multi-resolver v1.0.1
	github.com/fatih/color v1.13.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/go-hclog v1.1.0
	github.com/hashicorp/raft v1.3.3
	github.com/hashicorp/raft-boltdb v0.0.0-20211202195631-7d34b9fb3f42
	github.com/jessevdk/go-flags v1.5.0
	go.uber.org/zap v1.20.0
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	moul.io/number-to-words v0.6.0
)

require (
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20220114011407-0dd24b26b47d // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220114231437-d2e6a121cae0 // indirect
)

// replace github.com/Jille/grpc-multi-resolver => /home/developer/projects/go/grpc-multi-resolver
