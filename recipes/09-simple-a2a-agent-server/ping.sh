#!/bin/bash
: <<'COMMENT'
# Card discovery
COMMENT

HTTP_PORT=8080
AGENT_BASE_URL=http://0.0.0.0:${HTTP_PORT}

# host.docker.internal

curl ${AGENT_BASE_URL}/.well-known/agent.json \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  | jq '.'


