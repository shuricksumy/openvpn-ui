#!/bin/bash

set -a
source .env

docker-compose --env-file .env -f docker-compose.yml up