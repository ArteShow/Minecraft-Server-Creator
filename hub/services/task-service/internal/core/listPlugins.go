package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ListPlugins(serverID, token string) ([]string, error) {
	targetIP, err := ResolveServerTarget(serverID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "http://"+targetIP+":8003/server-service/plugin/list?server_id="+serverID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) > 0 {
			return nil, fmt.Errorf("host plugin list failed, status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("host plugin list failed, status %d", resp.StatusCode)
	}

	var result struct {
		Plugins []string `json:"plugins"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Plugins, nil
}
