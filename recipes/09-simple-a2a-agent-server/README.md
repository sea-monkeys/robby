
**Discovering**:
```bash
# This should return the agent metadata
curl -X GET http://0.0.0.0:8080/.well-known/agent.json \
  -H "Content-Type: application/json" \
  | jq '.'
```

c



**ask_for_something**:
```bash
# This should echo back the message 
curl -X POST http://0.0.0.0:8080 \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "id": "1111",
    "method": "message/send",
    "params": {
      "message": {
        "role": "user",
        "parts": [
          {
            "text": "What is the best pizza in the world?"
          }
        ]
      },
      "metadata": {
        "skill": "ask_for_something"
      }
    }
  }' \
  | jq '.'
```




**say_hello_world**:
```bash
curl -X POST http://0.0.0.0:8080 \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
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
  }' \
  | jq '.'
```
> NOTE: Streaming the answer could be insteresting


**another_task**:
```bash
curl -X POST http://0.0.0.0:8080 \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "id": "2222",
    "method": "message/send",
    "params": {
      "message": {
        "role": "user",
        "parts": [
          {
            "text": "Why the sky is blue?"
          }
        ]
      },
      "metadata": {
        "skill": "another_task"
      }
    }
  }' \
  | jq '.'
```

