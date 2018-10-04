#!/usr/bin/env bash

. ./.circleci/export-env-vars.sh

eval $(aws ecr get-login --no-include-email --region ap-south-1)

docker push ${IMAGE_NAME}:${IMAGE_TAG}

docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:${ENV}

docker push ${IMAGE_NAME}:${ENV}