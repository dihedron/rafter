package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	"github.com/dihedron/rafter/logging"
	pb "github.com/dihedron/rafter/proto"
	"github.com/hashicorp/raft"
)

type RPCInterface struct {
	pb.UnimplementedLogServer
	cache  *Cache
	raft   *raft.Raft
	logger logging.Logger
}

func NewRPCInterface(c *Cache, r *raft.Raft, l logging.Logger) *RPCInterface {
	return &RPCInterface{
		cache:  c,
		raft:   r,
		logger: l,
	}
}

func (r RPCInterface) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	message := &Message{
		Type: Get,
		Key:  request.Key,
	}
	data, err := json.Marshal(message)
	if err != nil {
		r.logger.Error("error marshalling Get message to JSON: %v", err)
		return nil, err
	}

	f := r.raft.Apply(data, time.Second)
	if err := f.Error(); err != nil {
		r.logger.Error("error masrhalling Get message to JSON: %v", err)
		return nil, rafterrors.MarkRetriable(err)
	}

	message = f.Response().(*Message)

	return &pb.GetResponse{
		Key:   request.Key,
		Value: message.Value,
	}, nil
	// return &pb.AddWordResponse{
	// 	CommitIndex: f.Index(),
	// }, nil
}

// func (r RPCInterface) AddWord(ctx context.Context, req *pb.AddWordRequest) (*pb.AddWordResponse, error) {
// 	f := r.raft.Apply([]byte(req.GetWord()), time.Second)
// 	if err := f.Error(); err != nil {
// 		return nil, rafterrors.MarkRetriable(err)
// 	}
// 	return &pb.AddWordResponse{
// 		CommitIndex: f.Index(),
// 	}, nil
// }

// func (r RPCInterface) GetWords(ctx context.Context, req *pb.GetWordsRequest) (*pb.GetWordsResponse, error) {
// 	r.wordTracker.mtx.RLock()
// 	defer r.wordTracker.mtx.RUnlock()
// 	return &pb.GetWordsResponse{
// 		BestWords:   cloneWords(r.wordTracker.words),
// 		ReadAtIndex: r.raft.AppliedIndex(),
// 	}, nil
// }

// func (r RPCInterface) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
// 	f := r.raft.Apply([]byte(req.Set()), time.Second)
// 	if err := f.Error(); err != nil {
// 		return nil, rafterrors.MarkRetriable(err)
// 	}
// 	return &pb.AddWordResponse{
// 		CommitIndex: f.Index(),
// 	}, nil
// }
