module github.com/terraform-providers/terraform-provider-vault

go 1.12

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/aws/aws-sdk-go v1.44.191
	github.com/gosimple/slug v1.4.1
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/terraform-plugin-sdk v1.1.1
	github.com/hashicorp/vault v1.13.5
	github.com/hashicorp/vault/api v1.9.0
	github.com/hashicorp/vault/sdk v0.8.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rainycape/unidecode v0.0.0-20150907023854-cb7f23ec59be // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
