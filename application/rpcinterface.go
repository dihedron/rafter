package application

import (
	"context"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	pb "github.com/dihedron/rafter/proto"
	"github.com/hashicorp/raft"
)

type RPCInterface struct {
	pb.UnimplementedExampleServer
	wordTracker *WordTracker
	raft        *raft.Raft
}

func NewRPCInterface(wt *WordTracker, r *raft.Raft) *RPCInterface {
	return &RPCInterface{
		wordTracker: wt,
		raft:        r,
	}
}

func (r RPCInterface) AddWord(ctx context.Context, req *pb.AddWordRequest) (*pb.AddWordResponse, error) {
	f := r.raft.Apply([]byte(req.GetWord()), time.Second)
	if err := f.Error(); err != nil {
		return nil, rafterrors.MarkRetriable(err)
	}
	return &pb.AddWordResponse{
		CommitIndex: f.Index(),
	}, nil
}

func (r RPCInterface) GetWords(ctx context.Context, req *pb.GetWordsRequest) (*pb.GetWordsResponse, error) {
	r.wordTracker.mtx.RLock()
	defer r.wordTracker.mtx.RUnlock()
	return &pb.GetWordsResponse{
		BestWords:   cloneWords(r.wordTracker.words),
		ReadAtIndex: r.raft.AppliedIndex(),
	}, nil
}

// func (r RPCInterface) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
// 	f := r.raft.Apply([]byte(req.Set()), time.Second)
// 	if err := f.Error(); err != nil {
// 		return nil, rafterrors.MarkRetriable(err)
// 	}
// 	return &pb.AddWordResponse{
// 		CommitIndex: f.Index(),
// 	}, nil
// }
