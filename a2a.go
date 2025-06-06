package robby

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var taskMutex sync.Mutex // Global mutex for protecting concurrent access

// IMPORTANT: minimal implementation of an A2A server for Robby Agents

// Serve the Agent Card at the well-known URL
func (agent *Agent) getAgentCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent.AgentCard)
}


// Alternative synchronous implementation that should work better
func (agent *Agent) handleTaskSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskRequest TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}

	switch taskRequest.Method {
	case "message/send":
		if len(taskRequest.Params.Message.Parts) > 0 {
			// Process the task synchronously without mutex in the HTTP handler
			// The mutex should only be in the AgentCallback if needed
			responseTask, err := agent.AgentCallback(taskRequest)
			if err != nil {
				log.Printf("Agent callback failed for task %s: %v", taskRequest.ID, err)
				http.Error(w, `{"error": "agent callback failed"}`, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(responseTask)
		} else {
			http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, `{"error": "unknown method"}`, http.StatusBadRequest)
	}
}

func (agent *Agent) StartA2AServer(addr string) error {
	logger := log.New(os.Stderr, "[Robby A2A] ", log.LstdFlags|log.Lshortfile)
	logger.Printf("Starting A2A Server on %s", addr)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/.well-known/agent.json", agent.getAgentCard)

	// Use the synchronous handler instead
	mux.HandleFunc("/", agent.handleTaskSync)

	// Start the server
	err := http.ListenAndServe(addr, mux)
	return err
}

func (agent *Agent) Ping(agentBaseURL string) (AgentCard, error) {
	resp, err := http.Get(agentBaseURL + "/.well-known/agent.json")
	if err != nil {
		return AgentCard{}, err
	}
	defer resp.Body.Close()

	var agentCard AgentCard
	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&agentCard); err != nil {
			return AgentCard{}, err
		}
		return agentCard, nil
	} else {
		return agentCard, errors.New("failed to ping agent: " + resp.Status)
	}
}

func (agent *Agent) SendToAgent(agentBaseURL string, taskRequest TaskRequest) (TaskResponse, error) {
	jsonTaskRequest, err := TaskRequestToJSONString(taskRequest)
	if err != nil {
		return TaskResponse{}, err
	}

	resp, err := http.Post(agentBaseURL+"/", "application/json", strings.NewReader(jsonTaskRequest))
	if err != nil {
		return TaskResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TaskResponse{}, errors.New("failed to send task request: " + resp.Status)
	}

	var taskResponse TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResponse); err != nil {
		return TaskResponse{}, err
	}

	return taskResponse, nil
}
