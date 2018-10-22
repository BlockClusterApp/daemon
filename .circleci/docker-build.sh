#!/usr/bin/env bash
. ./.circleci/export-env-vars.sh

docker build -f docker/Dockerfile \
    -t $IMAGE_NAME \
    .

docker tag $IMAGE_NAME 402432300121.dkr.ecr.us-west-2.amazonaws.com/blockcluster-daemon:latest