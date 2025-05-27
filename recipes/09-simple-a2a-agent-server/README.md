
**Discovering**:
```bash
# This should return the agent metadata
curl -X GET http://0.0.0.0:8080/.well-known/agent.json \
  -H "Content-Type: application/json" \
  | jq '.'
```

**ask_for_something**:
```bash
# This should echo back the message 
curl -X POST http://0.0.0.0:8080/tasks/send \
  -H "Content-Type: application/json" \
  -d '{
    "id": "ask_for_something",
    "message": {
      "role": "user",
      "parts": [
        {
          "text": "What is the best pizza in the world?"
        }
      ]
    }
  }' \
  | jq '.'
```

**say_hello_world**:
```bash
curl -X POST http://0.0.0.0:8080/tasks/send \
  -H "Content-Type: application/json" \
  -d '{
    "id": "say_hello_world",
    "message": {
      "role": "user", 
      "parts": [
        {
          "text": "Bob Morane"
        }
      ]
    }
  }' \
  | jq '.'
```
> NOTE: Streaming the answer could be insteresting


**another_task**:
```bash
curl -X POST http://0.0.0.0:8080/tasks/send \
  -H "Content-Type: application/json" \
  -d '{
    "id": "another_task",
    "message": {
      "role": "user", 
      "parts": [
        {
          "text": "Why the sky is blue?"
        }
      ]
    }
  }' \
  | jq '.'
```



```bash
# Test error handling - malformed request (missing parts)
curl -X POST http://0.0.0.0:8080/tasks/send \
  -H "Content-Type: application/json" \
  -d '{
    "id": "task-error",
    "message": {
      "role": "user"
    }
  }' \
  | jq '.'
```


```bash
#Test method not allowed (GET on /tasks/send)
curl -X GET http://0.0.0.0:8080/tasks/send \
  -H "Content-Type: application/json"
```