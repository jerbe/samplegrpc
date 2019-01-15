package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"os"
	"regexp"
	"github.com/jerbe/samplegrpc/proto/im"
	"strconv"
	"strings"
)

/**
    @author : Jerbe - The porter from Earth
    @time : 2019/1/14 12:00 PM
    @describe : 
*/

const ListenAddress = "0.0.0.0:9906"

type Client struct {

}

//简单RPC
func SampleCall(client im.IMClient,inputReader *bufio.Reader){
	fmt.Println("[Sample] please input message")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
	}

	response, err := client.Sample(context.Background(), &im.Request{Message:input})
	if err != nil {
		fmt.Println("[Sample] error:",err.Error())
	}else{
		fmt.Println("[Sample] response:", response.Message)
	}
}

//请求流
func RequestStreamCall(client im.IMClient,inputReader *bufio.Reader){
	fmt.Println("[RequestStream] please input message")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	requestClient, err := client.RequestStream(context.Background())
	if err != nil {
		fmt.Println("[RequestStream] error:",err.Error())
		return
	}

	re,_ := regexp.Compile(`^(-n (\d+))\s?`)
	matchStrings := re.FindStringSubmatch(input)
	repeatCount := 1
	if len(matchStrings) > 1 {
		repeatCount,_ = strconv.Atoi(matchStrings[2])
		input = re.ReplaceAllString(input,"")
	}

	for i:=1; i <= repeatCount; i++{
		err = requestClient.Send(&im.Request{Message:input+fmt.Sprintf("%d",i)})
		if err != nil {
			fmt.Println("[RequestStream] error:",err.Error())
		}
	}

	response,err := requestClient.CloseAndRecv()
	if err != nil {
		fmt.Println("[RequestStream] error:",err.Error())
	}

	fmt.Println("[RequestStream]", response.Message)
}

//返回流
func ResponseStreamCall(client im.IMClient,inputReader *bufio.Reader){
	fmt.Println("[RequestStream] please input message")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//调用ResponseStream方法
	responseStreamClient, err := client.ResponseStream(context.Background(), &im.Request{Message:input})
	if err != nil {
		fmt.Println("[RequestStream] error:",err.Error())
	}else{
		for{
			response,err := responseStreamClient.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("[RequestStream] error:",err.Error())
				continue
			}
			fmt.Println("[RequestStream] ", response.Message)
		}
	}
}

//双向流
func BilateralStreamCall(client im.IMClient,inputReader *bufio.Reader){
	fmt.Println("[BilateralStream] please input message")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bilateralStreamClient, err := client.BilateralStream(context.Background())
	if err != nil {
		fmt.Println("[BilateralStream] error:",err.Error())
		return
	}

	re,_ := regexp.Compile(`^(-n (\d+))\s?`)
	matchStrings := re.FindStringSubmatch(input)
	repeatCount := 1
	if len(matchStrings) > 1 {
		repeatCount,_ = strconv.Atoi(matchStrings[2])
		input = re.ReplaceAllString(input,"")
	}

	for i:=1; i <= repeatCount; i++{
		err = bilateralStreamClient.Send(&im.Request{Message:input+fmt.Sprintf("%d",i)})
		if err != nil {
			fmt.Println("[BilateralStream] send error:",err.Error())
		}
		response,err := bilateralStreamClient.Recv()
		if err != nil {
			fmt.Println("[BilateralStream] recv error:",err.Error())
		}else{
			fmt.Println("[BilateralStream] ", response.Message)
		}
	}

	err = bilateralStreamClient.CloseSend()
	if err != nil{
		fmt.Println("[BilateralStream] close send error:", err.Error())
	}

	for{
		response,err := bilateralStreamClient.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			fmt.Println("[BilateralStream] recv error:",err.Error())
			break
		}else{
			fmt.Println("[BilateralStream] ", response.Message)
		}
	}
}

func main(){
	conn,err := grpc.Dial(ListenAddress, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()

	client := im.NewIMClient(conn)
	inputReader := bufio.NewReader(os.Stdin)
	for{
		LoopStart:
		fmt.Println(`1: SampleRPCCall
2: RequestStreamCall
3: ResponseStreamCall
4: BilateralStreamCall
0: Exit
Please chose method:`)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		input = strings.Trim(input,"\n")
		fmt.Println(input)
		switch input {
		case "0":
			fmt.Println("System stopping...")
			os.Exit(0)
			case "1":
				SampleCall(client,inputReader)
			case "2":
				RequestStreamCall(client,inputReader)
			case "3":
				ResponseStreamCall(client,inputReader)
			case "4":
				BilateralStreamCall(client,inputReader)
			default:
				goto LoopStart
		}

	}
}
