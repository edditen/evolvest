#!/bin/bash

docker rmi -f $(docker images -a | egrep "evolvestd|<none>" | awk '{print $3}')
