module github.com/openshift-online/ocm-common

go 1.21

require (
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/aws-sdk-go-v2 v1.26.1
	github.com/aws/aws-sdk-go-v2/config v1.27.9
	github.com/aws/aws-sdk-go-v2/credentials v1.17.9
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.48.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.152.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.24.3
	github.com/aws/aws-sdk-go-v2/service/iam v1.27.1
	github.com/aws/aws-sdk-go-v2/service/kms v1.30.0
	github.com/aws/aws-sdk-go-v2/service/route53 v1.40.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.5
	github.com/hashicorp/go-version v1.6.0
	github.com/onsi/ginkgo/v2 v2.17.1
	github.com/onsi/gomega v1.30.0
	github.com/openshift-online/ocm-sdk-go v0.1.421
	github.com/sirupsen/logrus v1.9.3
	go.uber.org/mock v0.3.0
	golang.org/x/crypto v0.22.0
	gopkg.in/square/go-jose.v2 v2.6.0
)

require github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2 // indirect

require (
	github.com/aws/smithy-go v1.20.2
	github.com/kr/pretty v0.1.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.35.1
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.3 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
