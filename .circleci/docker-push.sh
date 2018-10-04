#!/usr/bin/env bash

. ./.circleci/export-env-vars.sh

eval $(aws ecr get-login --no-include-email --region us-west-2)

docker push ${IMAGE_NAME}:${IMAGE_TAG}

dokcer tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:${ENV}

docker push ${IMAGE_NAME}:${ENV}