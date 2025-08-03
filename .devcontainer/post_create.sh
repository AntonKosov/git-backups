#!/bin/bash

set -e

sudo chown devuser:devuser /workspaces
sudo chmod -R a+rw /go

mkdir -p $TMP_DIR

go install tool