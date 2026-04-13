package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PowerUsageResponse struct {
	Online     bool    `json:"online"`
	CPUPercent float64 `json:"cpu_percent"`
	RAMUsedMB  float64 `json:"ram_used_mb"`
	RAMLimitMB float64 `json:"ram_limit_mb"`
	RAMPercent float64 `json:"ram_percent"`
}

func GetPowerUsage(serverID, token, ownerID string) (PowerUsageResponse, error) {
	targetIP, err := ResolveServerTarget(serverID)
	if err != nil {
		return PowerUsageResponse{}, err
	}

	req, err := http.NewRequest("GET", "http://"+targetIP+":8003/server-service/power/usage?server_id="+serverID, nil)
	if err != nil {
		return PowerUsageResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Owner-ID", ownerID)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return PowerUsageResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) > 0 {
			return PowerUsageResponse{}, fmt.Errorf("host power usage failed, status %d: %s", resp.StatusCode, string(respBody))
		}
		return PowerUsageResponse{}, fmt.Errorf("host power usage failed, status %d", resp.StatusCode)
	}

	var result PowerUsageResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PowerUsageResponse{}, err
	}

	return result, nil
}
