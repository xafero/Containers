#!/bin/sh
echo go run imgview/src/imgview/program.go tcp://192.168.178.47:2375
go run imgfind/src/imgfind/program.go tcp://192.168.178.47:2375
