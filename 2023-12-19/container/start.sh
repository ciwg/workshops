#!/bin/bash

IMAGE="grok_chat"
CONTAINER="gpt_4_container"
DOCKERFILE="."

docker build -t ${IMAGE} -f ${DOCKERFILE} .

docker run --name ${CONTAINER} -d -p 8080:8080 ${IMAGE}
