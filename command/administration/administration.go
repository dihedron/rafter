// Binary raftadmin is a CLI interface to the RaftAdmin gRPC service.
package administration

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	proto "github.com/Jille/raftadmin/proto"
	"github.com/iancoleman/strcase"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/prototext"

	// Allow dialing multiple nodes with multi:///.
	_ "github.com/dihedron/grpc-multi-resolver"
	"github.com/dihedron/rafter/cluster"
	"github.com/dihedron/rafter/command/base"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	// Register health checker with gRPC.
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/health"
)

type Administration struct {
	base.Base

	Leader             bool           `short:"l" long:"leader" description:"Whether to dial the leader." optional:"yes"`
	Peers              []cluster.Peer `short:"t" long:"target" description:"The address of the node in the cluster to send the request to." required:"yes"`
	HealthCheckService string         `short:"h" long:"health-check" description:"Which gRPC service to health check when searching for the leader." optional:"yes" default:"quis.RaftLeader"`
}

func (cmd *Administration) Execute(args []string) error {

	logger := cmd.GetLogger(nil)
	ctx := context.Background()
	methods := proto.File_raftadmin_proto.Services().ByName("RaftAdmin").Methods()

	if len(args) < 1 {
		logger.Error("invalid command line format")
		var commands []string
		for i := 0; methods.Len() > i; i++ {
			commands = append(commands, strcase.ToKebab(string(methods.Get(i).Name())))
		}
		sort.Strings(commands)
		return fmt.Errorf("no command provided, please choose one of: %s", strings.Join(commands, ", "))
	}

	request, m, err := GetMessageByName(args[0], args[1:])
	if err != nil {
		return err
	}

	var serviceConfig grpc.DialOption = grpc.EmptyDialOption{}
	if cmd.Leader {
		c := fmt.Sprintf(`{"healthCheckConfig": {"serviceName": "%s"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`, cmd.HealthCheckService)
		serviceConfig = grpc.WithDefaultServiceConfig(c)
		logger.Debug("using service configuration: '%s'", c)
	}
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
	connection, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		serviceConfig,
	)

	// peers := []string{}
	// for _, peer := range cmd.Peers {
	// 	peers = append(peers, peer.Address.String())
	// }
	// address := fmt.Sprintf("multi:///%s", strings.Join(peers, ","))
	// logger.Info("connecting to %s", address)
	// retryOpts := []grpc_retry.CallOption{
	// 	grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
	// 	grpc_retry.WithMax(5),
	// }
	// connection, err := grpc.Dial(address,
	// 	grpc.WithDefaultServiceConfig(serviceConfig), grpc.WithInsecure(),
	// 	grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	// 	grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	// 	grpc.WithBlock())
	if err != nil {
		logger.Error("error connecting to gRPC server(s) %s: %v", address, err)
		return err
	}
	defer connection.Close()

	logger.Debug("invoking %s(%s)", m.Name(), prototext.Format(request.Interface()))
	response := messageFromDescriptor(m.Output()).Interface()
	if err := connection.Invoke(ctx, "/RaftAdmin/"+string(m.Name()), request.Interface(), response); err != nil {
		logger.Error("error invoking %s(%s): %v", m.Name(), prototext.Format(request.Interface()), err)
		return err
	}
	logger.Debug("response: '%s'", strings.TrimSpace(prototext.Format(response)))

	// this method returned a future; we call Await() to get the result,
	// and then Forget() to free up the memory of the server
	if f, ok := response.(*proto.Future); ok {
		c := proto.NewRaftAdminClient(connection)
		logger.Debug("invoking Await(%s)", strings.TrimRight(prototext.Format(f), "\n\r"))
		response, err := c.Await(ctx, f)
		if err != nil {
			logger.Error("error invoking Await(%s): %v", strings.TrimRight(prototext.Format(f), "\n\r"), err)
			return err
		}
		logger.Info("response: '%s'", strings.TrimSpace(prototext.Format(response)))
		fmt.Printf("%s(%s) --> %s\n", m.Name(), prototext.Format(request.Interface()), strings.TrimSpace(prototext.Format(response)))
		if _, err := c.Forget(ctx, f); err != nil {
			logger.Error("error invoking Forget(): %v", err)
			return err
		}
	} else {
		fmt.Printf("%s(%s) --> %s\n", m.Name(), prototext.Format(request.Interface()), strings.TrimSpace(prototext.Format(response)))
	}
	return nil
}
