#!/bin/bash -e
#
#

test_run_dir="$MINT_RUN_CORE_DIR/aws-sdk-go"
(cd "$test_run_dir" && GO111MODULE=on CGO_ENABLED=0 go build)
