#!/bin/bash

BASEDIR="/home/vagrant/go/src/github.wdf.sap.corp/velocity/axino"
LIBROOT="$BASEDIR/import/content/nwrfcsdk"
export LD_LIBRARY_PATH=$LIBROOT/lib
export CGO_LDFLAGS="-L$LIBROOT/lib -lsapnwrfc -lsapucum"
export CGO_CFLAGS="-I$LIBROOT/include"
export CGO_LDFLAGS_ALLOW=".*"
export CGO_CFLAGS_ALLOW=".*"
export DH_CONNECTION_URL=http://localhost:3000
export DH_CA_FOLDER=docker/ca
export DH_MOCKED_CONNECTION_DATA="/sapmnt/home/I335526/tmp/connection_data.json"
export DATA_INTEGRATION_SERVICE=true
export DH_CONNECTION_URL=http://localhost:3000
export DH_VSYSTEM_URL=http://localhost:8796

code