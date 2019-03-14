#!/bin/bash
set -e

#
# Script to take the generated CA cert and generate a k8s deployment manifest with it embedded
#

SCRIPT=$(readlink -f $0)
BASE_DIR=`dirname ${SCRIPT}`
GENERATED_DIR="$BASE_DIR/../generated"

export CABUNDLE=$(base64 $GENERATED_DIR/ca.crt -w 0)
envsubst  < webhook.yaml  > generated/webhook.yaml