#!/usr/bin/env bash

. ./.circleci/export-env-vars.sh

eval $(aws ecr get-login --no-include-email --region ap-south-1)

docker push ${IMAGE_NAME}
