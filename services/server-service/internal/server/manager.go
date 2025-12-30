package server

import (
	"os"

	eulaacceptor "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/eula_acceptor"
	getjar "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/get_jar"
	idgenerator "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/id_generator"
)

func CreateServer(version string) (string, error) {
	id := idgenerator.GenerateServerID()
	err := os.MkdirAll("./servers/"+id, 0755)
	if err != nil {
		return "", err
	}

	err = getjar.GetServerJar(version, "./server/"+id+"/")
	if err != nil {
		return "", err
	}

	err = eulaacceptor.WriteEULA("./server/" + id + "/")
	if err != nil {
		return "", err
	}

	return id, nil
}
