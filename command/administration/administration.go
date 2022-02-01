// Binary raftadmin is a CLI interface to the RaftAdmin gRPC service.
package administration

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	proto "github.com/Jille/raftadmin/proto"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/iancoleman/strcase"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/prototext"

	// Allow dialing multiple nodes with multi:///.
	_ "github.com/dihedron/grpc-multi-resolver"
	"github.com/dihedron/rafter/cluster"
	"github.com/dihedron/rafter/command/base"

	// Register health checker with gRPC.
	_ "google.golang.org/grpc/health"
)

type Administration struct {
	base.Base

	Leader             bool           `short:"l" long:"leader" description:"Whether to dial the leader." optional:"yes"`
	Peers              []cluster.Peer `short:"p" long:"peer" description:"The address of a peer node in the cluster to join" required:"yes"`
	HealthCheckService string         `short:"h" long:"health-check" description:"Which gRPC service to health check when searching for the leader." optional:"yes" default:"quis.RaftLeader"`
}

func (cmd *Administration) Execute(args []string) error {

	logger := cmd.GetLogger(nil)
	ctx := context.Background()
	methods := proto.File_raftadmin_proto.Services().ByName("RaftAdmin").Methods()

	if len(args) < 1 {
		logger.Error("invalid command line format")
		// if flag.NArg() < 2 {
		var commands []string
		for i := 0; methods.Len() > i; i++ {
			commands = append(commands, strcase.ToKebab(string(methods.Get(i).Name())))
		}
		sort.Strings(commands)
		return fmt.Errorf("no command provided, please choose one of: %s", strings.Join(commands, ", "))
	}

	req, m, err := GetMessageByName(args[0], args[1:])
	if err != nil {
		return err
	}
	serviceConfig := fmt.Sprintf(`{"healthCheckConfig": {"serviceName": "%s"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`, cmd.HealthCheckService)
	peers := []string{}
	for _, peer := range cmd.Peers {
		peers = append(peers, peer.Address.String())
	}
	address := fmt.Sprintf("multi:///%s", strings.Join(peers, ","))
	logger.Info("connecting to %s", address)
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}
	conn, err := grpc.Dial(address,
		grpc.WithDefaultServiceConfig(serviceConfig), grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Invoking %s(%s)", m.Name(), prototext.Format(req.Interface()))
	resp := messageFromDescriptor(m.Output()).Interface()
	if err := conn.Invoke(ctx, "/RaftAdmin/"+string(m.Name()), req.Interface(), resp); err != nil {
		return err
	}
	log.Printf("Response: %s", strings.TrimSpace(prototext.Format(resp)))

	// this method returned a future; we call Await() to get the result,
	// and then Forget() to free up the memory of the server
	if f, ok := resp.(*proto.Future); ok {
		c := proto.NewRaftAdminClient(conn)
		log.Printf("Invoking Await(%s)", strings.TrimRight(prototext.Format(f), "\n\r"))
		resp, err := c.Await(ctx, f)
		if err != nil {
			return err
		}
		log.Printf("Response: %s", prototext.Format(resp))
		if _, err := c.Forget(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
