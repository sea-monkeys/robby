#!/bin/bash
: <<'COMMENT'
# Send task to the agent server
COMMENT

HTTP_PORT=8080
AGENT_BASE_URL=http://0.0.0.0:${HTTP_PORT}

# host.docker.internal

read -r -d '' DATA <<- EOM
{
    "jsonrpc": "2.0",
    "id": "2222",
    "method": "message/send",
    "params": {
      "message": {
        "role": "user",
        "parts": [
          {
            "text": "Bob Morane"
          }
        ]
      },
      "metadata": {
        "skill": "say_hello_world"
      }
    }
}
EOM

curl ${AGENT_BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d "${DATA}" | jq '.'


