#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# necessary if the secrets were removed for some reason
# down would fail because env_file not found
(cd $SCRIPT_DIR && ./init_secrets.sh)

COMPOSE_PROJECT_NAME=core docker-compose -f deployments/webapp.docker-compose.yaml down --remove-orphans
