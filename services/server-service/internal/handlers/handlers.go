package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/models"
	getjar "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/get_jar"
	idgenerator "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/id_generator"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var req models.CreateServerRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := idgenerator.GenerateServerID()

	err = os.Mkdir("./servers"+id, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = getjar.GetServerJar(req.Version, "./servers/"+id+"/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var res models.CreateServerResponse
	res.ServerID = id

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
