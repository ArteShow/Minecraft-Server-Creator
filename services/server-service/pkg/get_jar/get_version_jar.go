package getjar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func GetServerJar(version, destDir string) error {
	fmt.Println("Creating server folder:", destDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	fmt.Println("Downloading version_manifest.json")
	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get version_manifest.json: %s", resp.Status)
	}

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
		return fmt.Errorf("version %s not found", version)
	}

	fmt.Println("Downloading version JSON for version:", version)
	resp2, err := http.Get(versionURL)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get version JSON: %s", resp2.Status)
	}

	var vjson VersionJSON
	if err := json.NewDecoder(resp2.Body).Decode(&vjson); err != nil {
		return err
	}

	jarURL := vjson.Downloads.Server.URL
	if jarURL == "" {
		return fmt.Errorf("server.jar URL not found for version %s", version)
	}

	fmt.Println("Downloading server.jar from:", jarURL)
	resp3, err := http.Get(jarURL)
	if err != nil {
		return err
	}
	defer resp3.Body.Close()

	if resp3.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download server.jar: %s", resp3.Status)
	}

	jarFile := filepath.Join(destDir, "server.jar")
	fmt.Println("Saving server.jar to:", jarFile)
	out, err := os.Create(jarFile)
	if err != nil {
		return err
	}
	defer out.Close()

	n, err := io.Copy(out, resp3.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded %d bytes\n", n)
	return nil
}
