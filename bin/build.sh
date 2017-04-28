#!/bin/sh

source `/usr/bin/dirname $0`/_check.sh

source `/usr/bin/dirname $0`/_setgopath.sh
echo "use GOPATH: $GOPATH"

echo "install dependencies..."
pushd `/usr/bin/dirname $0`/../src/go-tag-predict
GOPATH=$GOPATH glide install
popd

echo "build..."
cd `/usr/bin/dirname $0`/../
GOPATH=$GOPATH go build -o bin/tag-predict go-tag-predict
