package get_version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

type paperBuildsResponse struct {
	Builds []int `json:"builds"`
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
	for _, versionEntry := range manifest.Versions {
		if versionEntry.ID == version {
			versionURL = versionEntry.URL
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

	var versionJSON VersionJSON
	if err := json.NewDecoder(resp2.Body).Decode(&versionJSON); err != nil {
		return nil, err
	}

	jarURL := versionJSON.Downloads.Server.URL
	if jarURL == "" {
		return nil, errors.New("server.jar URL not found for version " + version)
	}

	return downloadAndCacheFile(jarURL, jarPath)
}

func GetPaperJar(version string) ([]byte, error) {
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return nil, err
	}

	jarPath := filepath.Join(cacheDir, "paper-"+version+".jar")
	if _, err := os.Stat(jarPath); err == nil {
		return os.ReadFile(jarPath)
	}

	buildsURL := "https://api.papermc.io/v2/projects/paper/versions/" + version + "/builds"
	resp, err := http.Get(buildsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paper builds API: %s", resp.Status)
	}

	var builds paperBuildsResponse
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return nil, err
	}
	if len(builds.Builds) == 0 {
		return nil, fmt.Errorf("no Paper builds found for version %s", version)
	}

	latestBuild := builds.Builds[len(builds.Builds)-1]
	buildStr := strconv.Itoa(latestBuild)
	jarName := "paper-" + version + "-" + buildStr + ".jar"
	downloadURL := "https://api.papermc.io/v2/projects/paper/versions/" + version + "/builds/" + buildStr + "/downloads/" + jarName

	return downloadAndCacheFile(downloadURL, jarPath)
}

func GetSpigotJar(version string) ([]byte, error) {
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return nil, err
	}

	jarPath := filepath.Join(cacheDir, "spigot-"+version+".jar")
	if _, err := os.Stat(jarPath); err == nil {
		return os.ReadFile(jarPath)
	}

	downloadURL := "https://download.getbukkit.org/spigot/spigot-" + version + ".jar"
	data, err := downloadAndCacheFile(downloadURL, jarPath)
	if err == nil {
		return data, nil
	}

	return GetPaperJar(version)
}

func downloadAndCacheFile(downloadURL, destinationPath string) ([]byte, error) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(destinationPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	data, err := io.ReadAll(io.TeeReader(resp.Body, out))
	if err != nil {
		return nil, err
	}

	return data, nil
}
