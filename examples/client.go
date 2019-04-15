package main

import (
	"context"
	"fmt"
	"github.com/thewinds/msgwrapper"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// ok
	CallHello("thewinds")
	// 400
	CallHello("t")
}

func CallHello(name string) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(msgwrapper.MsgWrapperClientInterceptor))
	if err != nil {
		log.Fatalln(err)
	}

	client := NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &HelloRequest{
		Name: name,
	})
	// 将 error 转化为 statusError
	statusError := new(msgwrapper.Error).From(err)
	if statusError != nil {
		fmt.Printf("code = %d; message = %s\n", statusError.Code(), statusError.Message())
		return
	}
	fmt.Println(resp.Result)

}
