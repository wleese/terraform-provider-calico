# Calico Terraform Provider

## Install
```
$ go get github.com/wleese/terraform-provider-calico
$ cd into terraform-provider-calico/vendor/github.com/projectcalico/libcalico-go
$ run glide install
```

## Usage

### Provider Configuration
provider.tf
```
provider "calico" {
  backend_type = "etcdv2"
  backend_etcd_authority = "192.168.56.20:2379"
}
```
### Host Endpoint
```
resource "calico_hostendpoint" "myendpoint" {
  name = "myendpoint"
  node = "my-endpoint-001"
  interface = "eth0"
  expected_ips = ["127.0.0.1"]
  profiles = ["endpointprofile"]
  labels = { endpointlabel = "myvalue" }
}
```

## Testing
The script test.sh will:
- download calicoctl and terraform
- build terraform-provider-calico
- spin up a container with etcd (docker-compose)
- pull tests out of testing/test_*
- do a terraform apply of the TF file
- use calicoctl to get the result
- compare it with the prestored results in the test_*.yaml file
