package data

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/dihedron/grpc-multi-resolver"
	"github.com/dihedron/rafter/command/data/random"
	proto "github.com/dihedron/rafter/distributed/proto"
	"github.com/dihedron/rafter/logging"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/montanaflynn/stats"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
)

type Benchmark struct {
	Base
	Key         string `short:"k" long:"key" description:"The key to use for the benckmark" optional:"yes" default:"_benchmark"`
	Iterations  int    `short:"i" long:"iterations" description:"The number of distinct values to set during the benchmark" optional:"yes" default:"1000"`
	Concurrency int    `short:"c" long:"concurrency" description:"The number of goroutines to run in parallel" optional:"yes" default:"10"`
	Length      int    `short:"l" long:"length" description:"The length of the random values to set in the benchmark" optional:"yes" default:"16"`
	Verbose     bool   `short:"v" long:"verbose" description:"Whether to produce verbose output" optional:"yes"`
}

// run with ./rafter data benchmark --peer=@tests/raft/node1.json --peer=@tests/raft/node2.json --peer=@tests/raft/node3.json --length=64 --iterations=1000 --concurrency=10

func (cmd *Benchmark) Execute(args []string) error {

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
	logger.Info("connecting to %s", address)
	conn, err := grpc.Dial(address,
		grpc.WithDefaultServiceConfig(serviceConfig), grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	if err != nil {
		logger.Error("dialing failed: %v", err)
		return err
	}
	defer conn.Close()
	c := proto.NewContextClient(conn)

	ch := generateWords(cmd.Iterations, cmd.Length)

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < cmd.Concurrency; i++ {
		wg.Add(1)
		go func(goroutine int) {
			defer wg.Done()
			ts := []time.Duration{}
			for w := range ch {
				start := time.Now()
				_, err := c.Set(context.Background(), &proto.SetRequest{Key: cmd.Key, Value: []byte(w)})
				elapsed := time.Since(start)
				ts = append(ts, elapsed)
				if err != nil {
					logger.Error("[%d] Set RPC failed: %v", goroutine, err)
					os.Exit(1)
				}
				if cmd.Verbose {
					if cmd.Length < 64 {
						logger.Info("[%d] Setting '%s' value to '%s' took %s", goroutine, cmd.Key, w, elapsed)
					} else {
						logger.Info("[%d] Setting '%s' value to '%s...' took %s", goroutine, cmd.Key, w[0:16], elapsed)
					}
				}
			}
			data := stats.LoadRawData(ts)
			mean, _ := data.Mean()
			stddev, _ := data.StandardDeviation()
			logger.Info("[%d] Stats: mean %s, std dev: %s", goroutine, time.Duration(mean), time.Duration(stddev))
		}(i)
	}
	wg.Wait()
	_, err = c.Get(context.Background(), &proto.GetRequest{Key: cmd.Key})
	if err != nil {
		logger.Error("Get RPC failed: %v", err)
		os.Exit(1)
	}
	_, err = c.Remove(context.Background(), &proto.RemoveRequest{Key: cmd.Key})
	if err != nil {
		logger.Error("Remove RPC failed: %v", err)
		os.Exit(1)
	}
	elapsed := time.Since(start)
	logger.Info("[FINAL] Benchmark run took %s", elapsed)
	cmd.ProfileMemory(logger)
	return nil
}

func generateWords(number int, length int) <-chan string {
	ch := make(chan string, 1)
	go func() {
		for i := 1; number > i; i++ {
			ch <- random.String(length)
		}
		close(ch)
	}()
	return ch
}
