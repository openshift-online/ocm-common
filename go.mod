module github.com/openshift-online/ocm-common

go 1.21

toolchain go1.21.7

replace github.com/openshift-online/ocm-sdk-go => /home/thomasmckay/go/src/github.com/openshift-online/ocm-sdk-go

require (
	github.com/aws/aws-sdk-go-v2 v1.22.2
	github.com/aws/aws-sdk-go-v2/service/iam v1.27.1
	github.com/golang/mock v1.6.0
	github.com/hashicorp/go-version v1.6.0
	github.com/onsi/ginkgo/v2 v2.11.0
	github.com/onsi/gomega v1.27.8
	github.com/openshift-online/ocm-sdk-go v0.1.391
	go.uber.org/mock v0.3.0
	golang.org/x/crypto v0.20.0
	gopkg.in/square/go-jose.v2 v2.6.0
)

require (
	github.com/aws/smithy-go v1.16.0
	github.com/kr/pretty v0.1.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.9.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
