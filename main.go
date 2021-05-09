package main

import (


	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	common "github.com/tinglanant/go_study_common"

	"go.mod/domain/repository"
	"go.mod/handler"
	"go.mod/proto/category"
	dataService "go.mod/domain/service"

	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"

	_ "github.com/jinzhu/gorm/dialects/mysql"

)

var QPS = 1000

func main() {
	//配置中心
	consulConfig,err := common.GetConsulConfig("127.0.0.1",8500,"/micro/config")
	if err !=nil {
		log.Error(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	//链路追踪
	t,io,err:=common.NewTracer("go.micro.service.category","localhost:6831")
	if err !=nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
		//这里设置地址和需要暴露的端口
		micro.Address("127.0.0.1:8082"),
		//添加consul 作为注册中心
		micro.Registry(consulRegistry),
		//绑定链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	//获取mysql配置,路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig,"mysql")

	//连接数据库
	db,err := gorm.Open("mysql",mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止复表
	db.SingularTable(true)


	//rp := repository.NewCategoryRepository(db)
	//rp.InitTable()
	// Initialise service
	service.Init()

	categoryDataService := dataService.NewCategoryDataService(repository.NewCategoryRepository(db))

	 err = category.RegisterCategoryHandler(service.Server(),&handler.Category{CategoryDataService:categoryDataService})
	 if  err != nil {
	 	log.Error(err)
	 }

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
