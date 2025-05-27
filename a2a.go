package robby

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var taskRequest TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
		http.Error(w, `{"error": "invalid request format"}`, http.StatusBadRequest)
		return
	}

	if len(taskRequest.Message.Parts) > 0 {

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

}

func (agent *Agent) StartA2AServer(addr string) error {

	//addr := agent.AgentCard.URL

	logger := log.New(os.Stderr, "[Robby A2A] ", log.LstdFlags|log.Lshortfile)
	logger.Printf("Starting HelloWorld A2A Server on %s", addr)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/.well-known/agent.json", agent.getAgentCard)
	mux.HandleFunc("/tasks/send", agent.handleTask)

	// Start the server
	err := http.ListenAndServe(addr, mux)

	return err
}
