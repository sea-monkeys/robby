package robby

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

// IMPORTANT: minimmal implementation of an A2A server for Robby Agents

// Serve the Agent Card at the well-known URL
func (agent *Agent) getAgentCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent.AgentCard)
}

// Handle incoming task requests at the A2A endpoint
func (agent *Agent) handleTask(w http.ResponseWriter, r *http.Request) {
	// TODO: TEST: taskRequest.JSONRpcVersion

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskRequest TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		//fmt.Println("🔴 Error decoding task request:", err)
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}
	/*
		fmt.Println("🟠 Task Request ID:", taskRequest.ID)
		jsonTaskRequest, err := TaskRequestToJSONString(taskRequest)
		if err != nil {
			fmt.Println("🔴 Error converting task request to JSON:", err)
		}
		fmt.Println("📝 Task Request:", jsonTaskRequest)
	*/

	switch taskRequest.Method {
	case "message/send":
		if len(taskRequest.Params.Message.Parts) > 0 {

			responseTask, err := agent.AgentCallback(taskRequest)

			if err != nil {
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

	//addr := agent.AgentCard.URL

	logger := log.New(os.Stderr, "[Robby A2A] ", log.LstdFlags|log.Lshortfile)
	logger.Printf("Starting HelloWorld A2A Server on %s", addr)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/.well-known/agent.json", agent.getAgentCard)

	mux.HandleFunc("/", agent.handleTask)

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
	// Cast the payload to an Agent Card
	var agentCard AgentCard

	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&agentCard); err != nil {
			return AgentCard{}, err
		} else {
			return agentCard, err
		}
	} else {
		return agentCard, errors.New("failed to ping agent: " + resp.Status)
	}
}

// QUESTION: Rename it to SendMessage? SendToAgent? SendToAgent?
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
