#!/usr/bin/env bash

echo "-> Building Datastore Application..."
cd ./datastore
if go build ; then
    echo "  - Built Successfully!"
else
    echo "  - Error during building Datastore"
    exit -1
fi

echo "-> Building Encryption-Server Application..."
cd ../encryption
if go build ; then
    echo "  - Built Successfully!"
else
    echo "  - Error during building Encruption-Server Application"
    exit -1
fi

echo "-> Building Yoti Command-line Application..."
cd ../yoti
if go build ; then
    echo "  - Built Successfully!"
else
    echo "  - Error during building Yoti"
    exit -1
fi

echo "!!! Build Completed !!!"