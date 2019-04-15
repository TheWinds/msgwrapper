msgwrapper gRPC 消息封装拦截器
---

## Why
msgwrapper用于解决message外部定义业务error的问题，
特别是在一些服务监控中会认为gRPC返回了error是一次调用失败，这种情况下只能把错误信息写到message中，
但这样又会让message的定义和解析变得繁琐。
msgwrapper 通过在原有message之外加一层response message的方式，在不侵入原有message的情况下，
实现了error message 和 error code 等信息的返回。

## Features
- 无需在所有的`Protobuffer message`都中定义`code` `message` 等重复字段，低侵入性
- 除了使用拦截器和用`msgwrapper`的方法来创建和解析错误外，无需额外代码

## Todo
- 错误详情中默认写入错误栈信息

## 如何使用

参考examples里面的例子

服务端 :
```go
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
```go
func main() {
	// ok
	CallHello("thewinds")
	// 400 the length of name must gather than 3
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
	statusError := new(msgwrapper.Error).From(err)
	if statusError != nil {
		fmt.Printf("code = %d; message = %s\n", statusError.Code(), statusError.Message())
		return
	}
	fmt.Println(resp.Result)

}
```