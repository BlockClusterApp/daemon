#!/usr/bin/env bash

. ./.circleci/export-env-vars.sh

eval $(aws ecr get-login --no-include-email --region us-west-2)

docker push ${IMAGE_NAME}

docker push 402432300121.dkr.ecr.us-west-2.amazonaws.com/blockcluster-daemon:latest