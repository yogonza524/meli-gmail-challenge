#!/usr/bin/env bash

set -e
set -o pipefail

# Delete go.mod if exists
rm -rf go.mod | true

echo "[1/5] Build volume"
docker volume create --name golang

echo "[2/5] Init module"
docker run --rm -v golang -v ${PWD}:/app -w /app golang go mod init github.com/yogonza524/meli-gmail-challenge

echo "[3/5] Module download"
docker run --rm -v golang -v ${PWD}:/app -w /app golang go mod download

echo "[4/5] Building..."
docker run \
  -v golang \
  -v ${PWD}:/app \
  -v ${PWD}/tmp:/go/pkg/mod \
  -v ${PWD}/packages:/usr/local/go/src/meli/domain \
  -w /app \
  -e CGO_ENABLED=0 \
  -e GOOS=linux \
  golang go build -a -installsuffix cgo -o /app/appBuilt .

echo "[5/5] Deleting workspace..."
rm -rf go.mod go.sum | true

echo "Sucess! Your challenge is ready!"