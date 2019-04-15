package main

import (
	"context"
	"github.com/thewinds/msgwrapper"
	"google.golang.org/grpc"
	"log"
	"net"
)

type HelloService struct {
}

func (*HelloService) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	if len(req.Name) < 3 {
		// 用 msgwrapper.NewError 创建错误
		return nil, msgwrapper.NewError(400, "the length of must gather than 3")
	}
	return &HelloResponse{
		Result: "hello," + req.Name,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	// 注册拦截器
	s := grpc.NewServer(grpc.UnaryInterceptor(msgwrapper.MsgWrapperInterceptor))
	RegisterHelloServiceServer(s, new(HelloService))
	log.Println("Start HelloService on :50051")
	err = s.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
}
