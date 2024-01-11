#! /bin/bash

conn=`grep -G "^CONN=" .env | cut -d= -f2`
pushd sql/schema
goose postgres $conn $1
popd