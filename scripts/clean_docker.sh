#!/bin/bash

docker rmi -f $(docker images -a | grep  "evolvestd" | awk '{print $3}')
