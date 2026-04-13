package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func DeletePlugin(serverID, token, fileName string) error {
	targetIP, err := ResolveServerTarget(serverID)
	if err != nil {
		return err
	}

	payload := []byte(`{"server_id":"` + serverID + `","file_name":"` + fileName + `"}`)
	req, err := http.NewRequest("POST", "http://"+targetIP+":8003/server-service/plugin/delete", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) > 0 {
			return fmt.Errorf("host plugin delete failed, status %d: %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("host plugin delete failed, status %d", resp.StatusCode)
	}

	return nil
}
