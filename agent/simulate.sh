#!/bin/bash

docker run -it --rm \
  --privileged \
  -v /sys/kernel/debug:/sys/kernel/debug:rw \
  -v /lib/modules:/lib/modules:ro \
  -v /usr/src:/usr/src:ro \
  -v /etc/localtime:/etc/localtime:ro \
  --pid=host \
  --runtime=nvidia \
  -e NVIDIA_VISIBLE_DEVICES=all \
  ubuntu:latest