#!/bin/bash

docker compose -f ../docker-compose.yml up -d cockroachdb1 cockroachdb2
docker compose exec cockroachdb1 ./cockroach init --insecure --host=cockroachdb1:26257
docker compose -f ../docker-compose.yml up -d migrate app client
