package main

import (
	"golang.org/x/net/context"
)

type ServeBiz struct {
	*Serve
}

func (s *ServeBiz) Check(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return nothing, nil
}

func (s *ServeBiz) Add(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return nothing, nil
}

func (s *ServeBiz) Test(ctx context.Context, nothing *Nothing) (*Nothing, error) {
	return nothing, nil
}

func (s *ServeBiz) mustEmbedUnimplementedBizServer() {
	//TODO implement me
	panic("implement me")
}
