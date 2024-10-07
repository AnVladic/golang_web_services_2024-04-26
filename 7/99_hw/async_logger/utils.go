package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"strings"
)

func isMethodAllowed(allowedMethods []string, method string) bool {
	subSplit := strings.Split(method, "/")
	for _, allowed := range allowedMethods {
		subAllowedSplit := strings.Split(allowed, "/")
		disallow := false
		for i, subAllowed := range subAllowedSplit {
			if subAllowed != "*" && subAllowed != subSplit[i] {
				disallow = true
				break
			}
		}
		if !disallow {
			return true
		}
	}
	return false
}

func CheckUserAccess(ctx context.Context, serve *Serve) (*Event, error) {
	methodName, ok := grpc.Method(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "Unknown method")
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata not found in context")
	}

	consumers := md["consumer"]
	if len(consumers) == 0 {
		return nil, status.Error(codes.Unauthenticated, "metadata not found in context")
	}
	consumer := consumers[0]
	accessMethods := serve.accessUserMap[consumer]
	if accessMethods == nil {
		return nil, status.Error(codes.Unauthenticated, "access method not found in context")
	}

	if !isMethodAllowed(accessMethods, methodName) {
		return nil, status.Error(
			codes.Unauthenticated, "access method not found for this consumer")
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "remove address not found")
	}

	return &Event{
		Consumer: consumer,
		Method:   methodName,
		Host:     p.Addr.String(),
	}, nil
}
