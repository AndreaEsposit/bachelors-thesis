#!/usr/bin/bash

libwasmer_url=$(curl -s https://api.github.com/repos/wasmerio/wasmer-nightly/releases/latest | jq --raw-output '.assets[] | select(.name == "wasmer-linux-amd64.tar.gz") | .browser_download_url')
mkdir -p libwasmer
curl -L $libwasmer_url > release.tar.gz
tar xzvf release.tar.gz -C libwasmer
# export CGO_CFLAGS="-I$(pwd)/libwasmer/include/"
# export CGO_LDFLAGS="-Wl,-rpath,$(pwd)/libwasmer/lib/ -L$(pwd)/libwasmer/lib/ -lwasmer_c_api"
# just test -tags custom_wasmer_runtime


export CGO_CFLAGS="-I$(pwd)/wasmer/packaged/include"
export CGO_LDFLAGS="-Wl,-rpath,$(pwd)/../wasmer/target/release/ -L/$(pwd)/../wasmer/target/release/ -lwasmer_c_api"
cd wasmer && go test -v -tags custom_wasmer_runtime