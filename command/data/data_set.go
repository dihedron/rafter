package data

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/dihedron/grpc-multi-resolver"
	pb "github.com/dihedron/rafter/application/proto"
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

	serviceConfig := `{"healthCheckConfig": {"serviceName": "Log"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`
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
	c := pb.NewContextClient(conn)
	response, err := c.Set(context.Background(), &pb.SetRequest{Key: cmd.Key, Value: []byte(cmd.Value)})
	if err != nil {
		log.Fatalf("Set RPC failed: %v", err)
		return err
	}
	fmt.Printf("key '%s' set to '%s' (index: %d)\n", cmd.Key, cmd.Value, response.Index)

	// var wg sync.WaitGroup
	// for i := 0; 10 > i; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		for w := range ch {
	// 			_, err := c.Set(context.Background(), &pb.SetRequest{Key: "key", Value: w})
	// 			if err != nil {
	// 				log.Fatalf("Set RPC failed: %v", err)
	// 			}
	// 		}
	// 	}()
	// }
	// wg.Wait()
	// resp, err := c.Get(context.Background(), &pb.GetRequest{Key: "key"})
	// if err != nil {
	// 	log.Fatalf("Get RPC failed: %v", err)
	// }
	// fmt.Println(resp)
	cmd.ProfileMemory(logger)
	return nil
}
