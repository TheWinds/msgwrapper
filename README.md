msgwrapper rpc 消息封装拦截器，用于解决message外部定义业务error的问题
---

## Features
- 无需在所有的`Protobuffer message`都中定义`code` `message` 等重复字段，低侵入性
- 除了使用拦截器和用msgwrapper的方法来创建和解析错误外，无需额外代码

## Todo
- 错误详情中默认写入错误栈信息

## 如何使用

参考examples里面的例子
```
├── Makefile
├── client.go
├── hello.pb.go
├── hello.proto
└── server.go
```

服务端 :
```
func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(msgwrapper.MsgWrapperInterceptor))
	hello.RegisterHelloServiceServer(s, new(HelloService))
	log.Println("Start HelloService on :50051")
	err = s.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
}
```

客户端:
```
func CallHello(name string) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(msgwrapper.MsgWrapperClientInterceptor))
	if err != nil {
		log.Fatalln(err)
	}

	client := NewHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &HelloRequest{
		Name: name,
	})
	statusError := new(msgwrapper.Error).From(err)
	if statusError != nil {
		fmt.Printf("code = %d; message = %s\n", statusError.Code(), statusError.Message())
		return
	}
	fmt.Println(resp.Result)

}
```