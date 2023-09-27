#!/bin/sh

protoc -I.:${GOPATH}/src  --gofast_out=plugins=ripple:. proto/*.proto
