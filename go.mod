module  github.com/tinglanant/go_study_category

go 1.15

require (
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/config/source/consul/v2 v2.9.1 // indirect
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2 v2.9.1 // indirect
	github.com/micro/go-plugins/wrapper/select/roundrobin/v2 v2.9.1 // indirect
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/common v0.23.0
	github.com/tinglanant/go_study_common v0.0.0-20210508094011-0727f16ba7c0
	google.golang.org/protobuf v1.26.0

)

replace github.com/tinglanant/go_study_common => ../go_study_common
