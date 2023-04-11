#!/bin/bash

#protoc --proto_path=proto --go_out=out --go_opt=paths=source_relative proto/album.proto

protoc --go_out=. --go_opt=paths=source_relative proto/*.proto