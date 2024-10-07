package main

import (
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

func authInterceptor(serve *Serve) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if err := serve.authInterceptor(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func logsInterceptor(serve *Serve) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			return resp, err
		}
		if err = serve.logsInterceptor(ctx); err != nil {
			return resp, err
		}
		return resp, nil
	}
}

func streamAuthInterceptor(serve *Serve) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if err := serve.authInterceptor(ss.Context()); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func streamLogInterceptor(serve *Serve) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		go func() {
			_ = serve.logsInterceptor(ss.Context())
		}()
		if err := handler(srv, ss); err != nil {
			return err
		}
		return nil
	}
}

func StartMyMicroservice(ctx context.Context, addr string, data string) error {
	var accessUserMap map[string][]string
	var logChannel = make(chan *Event)
	err := json.Unmarshal([]byte(data), &accessUserMap)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	serve := Serve{
		accessUserMap: accessUserMap,
		logChannel:    logChannel,
		listeners:     &[]Admin_LoggingServer{},
		stats:         make(map[Admin_StatisticsServer]*Stat),
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authInterceptor(&serve),
			logsInterceptor(&serve),
		),
		grpc.ChainStreamInterceptor(
			streamAuthInterceptor(&serve),
			streamLogInterceptor(&serve),
		),
	)

	RegisterAdminServer(server, &ServeAdmin{
		Serve: &serve,
	})
	RegisterBizServer(server, &ServeBiz{
		Serve: &serve,
	})

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				server.Stop()
				return
			case event := <-serve.logChannel:
				for _, server := range *serve.listeners {
					err = server.Send(event)
				}
			}
		}
	}()

	return nil
}
