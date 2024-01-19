#! /bin/bash

#conn=`grep -G "^CONN=" .env | cut -d= -f2`
source .env
pushd sql/schema
goose postgres $CONN $1
popd