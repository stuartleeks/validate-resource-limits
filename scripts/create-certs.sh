#!/bin/bash
set -e

#
# Script to generate a CA and the the certs for SSL in the webhook signed by the CA
# The CA is passed specified in the config for the admission controller
#

SCRIPT=$(readlink -f $0)
BASE_DIR=`dirname ${SCRIPT}`
GENERATED_DIR="$BASE_DIR/../generated"

# Clean old filess
rm $GENERATED_DIR -rf
mkdir $GENERATED_DIR

# Create Cert Authority
openssl genrsa -out $GENERATED_DIR/ca.key 2048
openssl req -new -x509 -key $GENERATED_DIR/ca.key -out $GENERATED_DIR/ca.crt -config $BASE_DIR/ca.config 

# Create Cert Signing Request and app cert signed with CA
openssl genrsa -out $GENERATED_DIR/app.key 2048
openssl req -new -key $GENERATED_DIR/app.key -subj "/CN=validateresourcelimits.default.svc" -out $GENERATED_DIR/app.csr -config $BASE_DIR/app.config
openssl x509 -req -in $GENERATED_DIR/app.csr -CA $GENERATED_DIR/ca.crt -CAkey $GENERATED_DIR/ca.key -CAcreateserial -out $GENERATED_DIR/app.crt
