package getjar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type VersionManifest struct {
	Versions []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"versions"`
}

type VersionJSON struct {
	Downloads struct {
		Server struct {
			URL string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}

func GetServerJar(version, dest string) error {
	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var manifest VersionManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return err
	}

	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == version {
			versionURL = v.URL
			break
		}
	}
	if versionURL == "" {
		return fmt.Errorf("version not found")
	}

	resp2, err := http.Get(versionURL)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	var vjson VersionJSON
	if err := json.NewDecoder(resp2.Body).Decode(&vjson); err != nil {
		return err
	}

	resp3, err := http.Get(vjson.Downloads.Server.URL)
	if err != nil {
		return err
	}
	defer resp3.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp3.Body)
	return err
}
