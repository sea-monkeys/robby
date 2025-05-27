package robby

import "encoding/json"

func TaskRequestToJSONString(taskRequest TaskRequest) (string, error) {
	jsonData, err := json.MarshalIndent(taskRequest, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
