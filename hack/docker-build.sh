#! /bin/sh
if [ "${IMAGE_TAG}" = "" ]; then
  IMAGE_TAG=latest
fi
docker build ${BUILD_OPT} -t ${IMAGE_NAME}:${IMAGE_TAG} .
