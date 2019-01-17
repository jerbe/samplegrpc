#!/usr/bin/env bash
docker build -t samplegrpc:1.0.0 -f ./Dockerfile --network=host --no-cache ./