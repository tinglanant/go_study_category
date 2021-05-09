package main

import (
	"context"
	"fmt"
	common "github.com/tinglanant/go_study_common"

	go_micro_service_category "github.com/tinglanant/go_study_category/proto/category"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"log"
)

func main()  {
	//注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})


	//链路追踪
	tracerClientName := "go.micro.service.category.client"
	t,io,err:=common.NewTracer(tracerClientName,"localhost:6831")
	if err !=nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name(tracerClientName),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"),
		//添加注册中心
		micro.Registry(consul),
		//绑定链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		)

	categoryService:=go_micro_service_category.NewCategoryService("go.micro.service.category",service.Client())

	productAdd := &go_micro_service_category.CategoryRequest{
		CategoryName:"11",
		CategoryLevel:1,
		CategoryParent:0,

	}
	response,err:=categoryService.CreateCategory(context.TODO(),productAdd)
	if err !=nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}
