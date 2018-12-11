#!/bin/bash

docker build -t lotto_dashboard .
docker create --name dash lotto_dashboard
docker cp dash:/go/bin/lotto_dashboard lotto_dashboard
docker rm dash
