package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/hanzug/goS/app/index_platform/analyzer"
	"github.com/hanzug/goS/app/mapreduce/master"
	"github.com/hanzug/goS/config"
	"github.com/hanzug/goS/idl/pb/mapreduce"
	"github.com/hanzug/goS/loading"
	"github.com/hanzug/goS/pkg/discovery"
	logs "github.com/hanzug/goS/pkg/logger"
)

const (
	MapreduceServerName = "mapreduce"
)

func main() {
	loading.Loading()
	analyzer.InitSeg()

	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewRegister(etcdAddress, logs.LogrusObj)
	defer etcdRegister.Stop()

	grpcAddress := config.Conf.Services[MapreduceServerName].Addr[0]
	node := discovery.Server{
		Name: config.Conf.Domain[MapreduceServerName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()

	mapreduce.RegisterMapReduceServiceServer(server, master.GetMapReduceSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err = etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
