#!/bin/bash

systemctl stop ap

docker build -f /home/core/share/docker/Dockerfile.base -t pottava/docker-webui:base .
