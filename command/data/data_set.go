package data

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/dihedron/grpc-multi-resolver"
	proto "github.com/dihedron/rafter/distributed/proto"
	"github.com/dihedron/rafter/logging"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
)

type Set struct {
	Base
	Key   string `short:"k" long:"key" description:"The key to set/replace" required:"yes"`
	Value string `short:"v" long:"value" description:"The value to set/replace" required:"yes"`
}

func (cmd *Set) Execute(args []string) error {

	logger := logging.NewConsoleLogger(logging.StdOut)
	defer cmd.ProfileCPU(logger).Close()

	serviceConfig := `{"healthCheckConfig": {"serviceName": "quis.RaftLeader"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}
	peers := []string{}
	for _, peer := range cmd.Peers {
		peers = append(peers, peer.Address.String())
	}
	address := fmt.Sprintf("multi:///%s", strings.Join(peers, ","))
	// conn, err := grpc.Dial("multi:///localhost:50051,localhost:50052,localhost:50053",
	conn, err := grpc.Dial(address,
		grpc.WithDefaultServiceConfig(serviceConfig), grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	if err != nil {
		log.Fatalf("dialing failed: %v", err)
		return err
	}
	defer conn.Close()
	c := proto.NewContextClient(conn)
	response, err := c.Set(context.Background(), &proto.SetRequest{Key: cmd.Key, Value: []byte(cmd.Value)})
	if err != nil {
		log.Fatalf("Set RPC failed: %v", err)
		return err
	}
	fmt.Printf("key '%s' set to '%s' (index: %d)\n", cmd.Key, cmd.Value, response.Index)
	cmd.ProfileMemory(logger)
	return nil
}
