package main

import (
	"github.com/Markaplay-Game-Hosting/GoEventBot/ui"
	"io/fs"
	"net/http"
)

func (app *application) spaHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(ui.DistDirFS, "index.html")
	if err != nil {
		app.logger.Error("Failed to load index.html", "error", err)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
