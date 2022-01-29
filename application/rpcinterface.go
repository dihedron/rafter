package application

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	proto "github.com/dihedron/rafter/application/proto"
	"github.com/dihedron/rafter/logging"
	"github.com/hashicorp/raft"
)

type RPCInterface struct {
	proto.UnimplementedStateServer
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

func (r RPCInterface) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetResponse, error) {
	message := &Message{
		Type: Get,
		Key:  request.Key,
	}

	r.logger.Debug("message received: %s", logging.ToJSON(message))
	data, err := json.Marshal(message)
	if err != nil {
		r.logger.Error("error marshalling Get message to JSON: %v", err)
		return nil, err
	}

	f := r.raft.Apply(data, time.Second)
	if err := f.Error(); err != nil {
		r.logger.Error("error applying Get message to cluster: %v", err)
		return nil, rafterrors.MarkRetriable(err)
	}
	if f.Response() != nil {
		switch response := f.Response().(type) {
		case error:
			r.logger.Error("received error from FSM: %v", response)
			return nil, rafterrors.MarkRetriable(response)
		case []byte:
			message := &Message{}
			if err := json.Unmarshal(response, message); err != nil {
				r.logger.Error("error unmarshalling response to Get message from cluster: %v", err)
				return nil, rafterrors.MarkRetriable(err)
			}
			return &proto.GetResponse{
				Key:   request.Key,
				Value: message.Value,
				Index: f.Index(),
			}, nil
		}
	}
	return nil, rafterrors.MarkRetriable(fmt.Errorf("nil response"))
}

func (r RPCInterface) Set(ctx context.Context, request *proto.SetRequest) (*proto.SetResponse, error) {
	message := &Message{
		Type:  Set,
		Key:   request.Key,
		Value: request.Value,
	}
	r.logger.Debug("message received: %s", logging.ToJSON(message))
	data, err := json.Marshal(message)
	if err != nil {
		r.logger.Error("error marshalling Set message to JSON: %v", err)
		return nil, err
	}

	f := r.raft.Apply(data, time.Second)
	if err := f.Error(); err != nil {
		r.logger.Error("error applying Set message to cluster: %v", err)
		return nil, rafterrors.MarkRetriable(err)
	}
	if f.Response() != nil {
		switch response := f.Response().(type) {
		case error:
			r.logger.Error("received error from FSM: %v", response)
			return nil, rafterrors.MarkRetriable(response)
		case []byte:
			message := &Message{}
			if err := json.Unmarshal(response, message); err != nil {
				r.logger.Error("error unmarshalling response to Set message from cluster: %v", err)
				return nil, rafterrors.MarkRetriable(err)
			}
			return &proto.SetResponse{
				Index: f.Index(),
			}, nil
		}
	}

	return nil, rafterrors.MarkRetriable(fmt.Errorf("nil response"))
}

func (r RPCInterface) Remove(ctx context.Context, request *proto.RemoveRequest) (*proto.RemoveResponse, error) {
	message := &Message{
		Type: Remove,
		Key:  request.Key,
	}
	r.logger.Debug("message received: %s", logging.ToJSON(message))
	data, err := json.Marshal(message)
	if err != nil {
		r.logger.Error("error marshalling Remove message to JSON: %v", err)
		return nil, err
	}

	f := r.raft.Apply(data, time.Second)
	if err := f.Error(); err != nil {
		r.logger.Error("error applying Set message to cluster: %v", err)
		return nil, rafterrors.MarkRetriable(err)
	}

	if f.Response() != nil {
		switch response := f.Response().(type) {
		case error:
			r.logger.Error("received error from FSM: %v", response)
			return nil, rafterrors.MarkRetriable(response)
		case []byte:
			message := &Message{}
			if err := json.Unmarshal(response, message); err != nil {
				r.logger.Error("error unmarshalling response to Remove message from cluster: %v", err)
				return nil, rafterrors.MarkRetriable(err)
			}
			return &proto.RemoveResponse{
				Key:   message.Key,
				Value: message.Value,
				Index: f.Index(),
			}, nil
		}
	}

	return nil, rafterrors.MarkRetriable(fmt.Errorf("nil response"))
}

func (r RPCInterface) List(ctx context.Context, request *proto.ListRequest) (*proto.ListResponse, error) {
	if request.Filter != "" {
		// perform sanity check on regexp before sending to FSM
		if _, err := regexp.Compile(request.Filter); err != nil {
			return nil, err
		}
	}
	message := &Message{
		Type:   List,
		Filter: request.Filter,
	}
	r.logger.Debug("message received: %s", logging.ToJSON(message))
	data, err := json.Marshal(message)
	if err != nil {
		r.logger.Error("error marshalling List message to JSON: %v", err)
		return nil, err
	}

	f := r.raft.Apply(data, time.Second)
	if err := f.Error(); err != nil {
		r.logger.Error("error applying List message to cluster: %v", err)
		return nil, rafterrors.MarkRetriable(err)
	}

	if f.Response() != nil {
		switch response := f.Response().(type) {
		case error:
			r.logger.Error("received error from FSM: %v", response)
			return nil, rafterrors.MarkRetriable(response)
		case []byte:
			message := &Message{}
			if err := json.Unmarshal(response, message); err != nil {
				r.logger.Error("error unmarshalling response to List message from cluster: %v", err)
				return nil, rafterrors.MarkRetriable(err)
			}
			return &proto.ListResponse{
				Keys:  message.Keys,
				Index: f.Index(),
			}, nil
		}
	}

	return nil, rafterrors.MarkRetriable(fmt.Errorf("nil response"))
}

func (r RPCInterface) Clear(ctx context.Context, request *proto.ClearRequest) (*proto.ClearResponse, error) {
	message := &Message{
		Type:   Clear,
		Filter: request.Filter,
	}
	r.logger.Debug("message received: %s", logging.ToJSON(message))
	data, err := json.Marshal(message)
	if err != nil {
		r.logger.Error("error marshalling Clear message to JSON: %v", err)
		return nil, err
	}

	f := r.raft.Apply(data, time.Second)
	if err := f.Error(); err != nil {
		r.logger.Error("error applying Clear message to cluster: %v", err)
		return nil, rafterrors.MarkRetriable(err)
	}

	if f.Response() != nil {
		switch response := f.Response().(type) {
		case error:
			r.logger.Error("received error from FSM: %v", response)
			return nil, rafterrors.MarkRetriable(response)
		case []byte:
			message := &Message{}
			if err := json.Unmarshal(response, message); err != nil {
				r.logger.Error("error unmarshalling response to Clear message from cluster: %v", err)
				return nil, rafterrors.MarkRetriable(err)
			}
			return &proto.ClearResponse{
				Index: f.Index(),
			}, nil
		}
	}

	return nil, rafterrors.MarkRetriable(fmt.Errorf("nil response"))
}
