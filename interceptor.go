package msgwrapper

import (
	"context"
	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

const iNTERCEPTOR_KEY = "_MsgWrapperInterceptor_wUw05cMW"

// MsgWrapperInterceptor
// 消息包装拦截器
// 如果使用此拦截器可以在Service实现方法中直接返回 `StatusError`
// 此拦截器会对`StatusError`进行特殊处理,最终gRPC服务返回的 code为 0 `OK`
// 如果Service实现方法中没有返回`StatusError`而是其他error，则最终gRPC服务返回的 code为 2 `Unknown`
func MsgWrapperInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if len(data.Get(iNTERCEPTOR_KEY)) != 1 {
		return handler(ctx, req)
	}

	response := &Response{Code: 200}
	resp, err := handler(ctx, req)
	if err != nil {
		// 判断错误类型
		if isStatusError(err) {
			statusError := err.(StatusError)
			response.Message = statusError.Message()
			response.Detail = statusError.Detail()
			response.Code = statusError.Code()
			return response, nil
		}
		// 不支持的错误类型直接返回
		return resp, err
	}
	// 包装response
	if message, ok := resp.(proto.Message); ok {
		any, err := ptypes.MarshalAny(message)
		if err != nil {
			response.Code = 500
			response.Message = "服务器错误"
			response.Detail = "MsgWrapperInterceptor: " + err.Error()
		}
		response.Data = any
	}
	return response, nil
}

// MsgWrapperClientInterceptor
// 消息包装客户端拦截器
// 对 MsgWrapperInterceptor 产生的 Response 进行接解包处理
// 将返回响应以及`StatusError`,如果要获取错误详情请讲`error`断言为`StatusError`
func MsgWrapperClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	newCtx := metadata.AppendToOutgoingContext(ctx, iNTERCEPTOR_KEY, "MsgWrapper")
	// 保留invoke方法需要返回的reply类型
	dataReply := reply
	// 替换为包装的reply类型
	reply = &Response{}
	// 执行rpc请求
	err := invoker(newCtx, method, req, reply, cc, opts...)
	// 如果返回了rpc系统的error,则直接进行返回
	if err != nil {
		return err
	}
	// 解包处理
	if resp, ok := reply.(*Response); ok {
		msg := dataReply.(proto.Message)
		if resp.Code != 200 {
			return NewErrorWithDetail(resp.Code, resp.Message, resp.Detail)
		}
		err = ptypes.UnmarshalAny(resp.Data, msg)
		if err != nil {
			return NewError(502, "类型转换失败")
		}
		*(&reply) = msg
	}
	return nil
}
