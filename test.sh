#!/bin/bash

WD=$(pwd)
WORKDIR=tmp
TFVERSION=0.7.11
TFARCH=linux_amd64
TFURL="https://releases.hashicorp.com/terraform/${TFVERSION}/terraform_${TFVERSION}_${TFARCH}.zip"
CALICOVERSION=v1.0.0-beta
CALICOURL="https://github.com/projectcalico/calico-containers/releases/download/${CALICOVERSION}/calicoctl"

if ! [[ -d $WORKDIR ]]; then
  mkdir $WORKDIR
fi
cd $WORKDIR

if ! [[ -e terraform ]]; then
  echo "Downloading Terraform"
  curl -s $TFURL -o terraform_${TFVERSION}_${TFARCH}.zip
  if [[ $? -ne 0 ]]; then
    echo "Failed to download terraform"
    exit 1
  fi
  unzip terraform_${TFVERSION}_${TFARCH}.zip
fi

if ! [[ -e calicoctl ]]; then
  echo "Downloading Calicoctl"
  curl -s -L $CALICOURL -o calicoctl
  if [[ $? -ne 0 ]]; then
    echo "Failed to download calicoctl"
    exit 1
  fi
  chmod +x calicoctl
fi

cd "$WD"

echo "Downloading GO dependencies"
go get -v
if [[ $? -ne 0 ]]; then
  echo "Failed to download all dependencies"
  exit 1
fi

echo "Building terraform-provider-calico"
go build -v
if [[ $? -ne 0 ]]; then
  echo "Failed to build terraform-provider-calico"
  exit 1
fi
cp terraform-provider-calico $WORKDIR
cp testing/* $WORKDIR
cd $WORKDIR

if ! grep "${WD}/terraform-provider-calico" ~/.terraformrc 2>&1 > /dev/null; then
  echo
  echo "You'll have to change your ~/.terraform.rc file to include this"
  echo "if you want to continue running these tests:"
  echo
  echo "providers {"
  echo "  calico = \"${WD}/terraform-provider-calico\""
  echo "}"
  exit 1
fi

# cleanup just in case
docker stop $(docker-compose ps -q) &>/dev/null
docker-compose kill &>/dev/null
docker-compose rm &>/dev/null

echo "Setting up ETCD"
docker-compose run -d etcd
ETCD_AUTHORITY="$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker-compose ps -q)):2379"

if [[ "$ETCD_AUTHORITY" == "" ]]; then
  echo "Failed to get ETCD endpoint"
  exit 1
fi

echo "Wait 5s until ETCD starts up"
sleep 5s

rm -rf test; mkdir test
sed "s/PLACEHOLDER/${ETCD_AUTHORITY}/" provider.tf > test/provider.tf

cp terraform test/
cp terraform-provider-calico test/
cp calicoctl test/

for i in hostendpoints workloadendpoints profiles ippools bgppeers policies; do
  tffile="${WD}/testing/test_${i}.tf"
  if [[ -e $tffile ]]; then
    echo "Testing ${i}"
    cp "$tffile" test/
    cd test
    ./terraform apply
    if [[ $? -ne 0 ]]; then
      echo "Failed to terraform apply (${tffile})"
      exit 1
    fi
    rm "test_${i}.tf"

    ETCD_AUTHORITY="$ETCD_AUTHORITY" ./calicoctl get $i -o yaml > test.yaml
    if [[ $? -ne 0 ]]; then
      echo "Failed to talk to ETCD at ${ETCD_AUTHORITY}"
      exit 1
    fi
    if ! diff test.yaml "${WD}/testing/test_${i}.yaml"; then
      echo "Expected ${i} yaml and that from testing/test_${i}.yaml do not match"
      exit 1
    else
      echo "${i} - OK"
    fi
  else
    echo "Don't have a test for ${i} - skipping"
  fi
done
