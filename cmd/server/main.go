package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"regexp"
	"github.com/jerbe/samplegrpc/proto/im"
	"strconv"
)

/**
    @author : Jerbe - The porter from Earth
    @time : 2019/1/14 12:00 PM
    @describe : 
*/

const ListenAddress = "0.0.0.0:9906"

type Server struct {

}

//简单RPC
func (s *Server) Sample(ctx context.Context, request *im.Request) (response *im.Response, err error){
	fmt.Println("from client use [Sample]:", request.Message)
	return &im.Response{Message:"success"},nil
}

//返回流
func (s *Server) ResponseStream(request *im.Request, responseServer im.IM_ResponseStreamServer) error{
	fmt.Println("from client use [ResponseStream]:", request.Message)
	for i:=1; i<= 10; i++{
		err := responseServer.Send(&im.Response{Message: fmt.Sprintf("success [%d]",i)})
		if err != nil {
			return err
		}
	}
	return nil
}

//请求流
func (s *Server) RequestStream(requestServer im.IM_RequestStreamServer) error{
	for{
		request, err := requestServer.Recv()
		if err == io.EOF {
			return requestServer.SendAndClose(&im.Response{Message:"success"})
		}
		if err != nil {
			fmt.Println("from client use [RequestStream] error:", err.Error())
			return err
		}
		if request == nil {
			return errors.New("request is nil")
		}
		fmt.Println("from client use [RequestStream]:", request.Message)
	}
}

//双向流
func (s *Server) BilateralStream(bilateralServer im.IM_BilateralStreamServer) error{
	for{
		request, err := bilateralServer.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println("[BilateralStream] error:", err.Error())
			return err
		}
		if request == nil {
			return errors.New("[BilateralStream] request is nil")
		}
		message := request.Message
		burstCount := 1
		re,_:=regexp.Compile(`^-b\s?(\d+)`)
		matchStrings := re.FindStringSubmatch(message)
		if len(matchStrings) > 1{
			burstCount,_ = strconv.Atoi(matchStrings[1])
		}
		fmt.Println("[BilateralStream] request:", message)
		for i:=1; i <= burstCount; i++ {
			response := &im.Response{Message:fmt.Sprintf("success %d", i)}
			err = bilateralServer.Send(response)
			if err != nil {
				fmt.Println("[BilateralStream] error:", err.Error())
				return err
			}
			fmt.Println("[BilateralStream] send success")
		}
	}
}

func main(){
	listen,err := net.Listen("tcp",ListenAddress)
	if err != nil {
		fmt.Println(err.Error())
	}

	grpcServer := grpc.NewServer()
	im.RegisterIMServer(grpcServer,&Server{})
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("error:",err.Error())
	}


}
