#!/bin/sh

source `/usr/bin/dirname $0`/_check.sh

source `/usr/bin/dirname $0`/_setgopath.sh
echo "use GOPATH: $GOPATH"

cd `/usr/bin/dirname $0`/../src/go-tag-predict
GOPATH=$GOPATH go run -race main.go
