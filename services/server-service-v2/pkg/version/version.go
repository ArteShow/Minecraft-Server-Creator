package get_version

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const cacheDir = ".cache/versions"

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

func GetServerJar(version string) ([]byte, error) {
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return nil, err
	}

	jarPath := filepath.Join(cacheDir, version+".jar")

	if _, err := os.Stat(jarPath); err == nil {
		return os.ReadFile(jarPath)
	}

	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get version_manifest.json: " + resp.Status)
	}

	var manifest VersionManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}

	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == version {
			versionURL = v.URL
			break
		}
	}
	if versionURL == "" {
		return nil, errors.New("version " + version + " not found")
	}

	resp2, err := http.Get(versionURL)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get version JSON: " + resp2.Status)
	}

	var vjson VersionJSON
	if err := json.NewDecoder(resp2.Body).Decode(&vjson); err != nil {
		return nil, err
	}

	jarURL := vjson.Downloads.Server.URL
	if jarURL == "" {
		return nil, errors.New("server.jar URL not found for version " + version)
	}

	resp3, err := http.Get(jarURL)
	if err != nil {
		return nil, err
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download server.jar: " + resp3.Status)
	}

	out, err := os.Create(jarPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	data, err := io.ReadAll(io.TeeReader(resp3.Body, out))
	if err != nil {
		return nil, err
	}

	return data, nil
}
