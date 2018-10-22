#!/usr/bin/env bash

export COMMIT_HASH=${CIRCLE_SHA1}
if [ "$CIRCLE_TAG" = "production" ] || [ "$CIRCLE_BRANCH" = "master" ];
then
  export ENV="production";
  export CLUSTER_PREFIX="production-ap-south-1b";
  export IMAGE_NAME="blockcluster/blockcluster-daemon"
# elif [ "$CIRCLE_TAG" = "staging" ] || [ "$CIRCLE_BRANCH"  = "staging" ];
# then
#   export NODE_ENV=staging
#   export CLUSTER_PREFIX="dev";
# elif [ "$CIRCLE_TAG" = "test" ] || [ "$CIRCLE_BRANCH" = "test" ] || [ "$IS_TEST" = "1" ];
# then
#   export NODE_ENV=test
#   export CLUSTER_PREFIX="dev";
elif [ "$CIRCLE_TAG" = "dev" ] ||  [ "$CIRCLE_BRANCH" = "dev" ];
then
  export ENV="dev";
  export CLUSTER_PREFIX="dev";
  export IMAGE_NAME="402432300121.dkr.ecr.us-west-2.amazonaws.com/blockcluster-daemon:${ENV}"
fi


