package app

import (
	"encoding/json"
	"fmt"
	"github.com/Ki4EH/go-bash/internal/logger"
	"github.com/gorilla/mux"
	"net/http"
)

// InfoAll returns all commands
func (a *App) InfoAll(w http.ResponseWriter, r *http.Request) {
	commands, err := a.AllCommands()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("error on get all commands ", err)
		return
	}
	jsonCommands, err := json.Marshal(commands)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("error on marshalling all commands ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonCommands)
}

// NewCommand creates a new command in the database
func (a *App) NewCommand(w http.ResponseWriter, r *http.Request) {
	var command *Table
	json.NewDecoder(r.Body).Decode(&command)
	if err := a.InsertCommand(command); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		logger.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RemoveCommand removes a command or multiple commands from the database
func (a *App) RemoveCommand(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query().Get("id")
	if err := a.Remove(params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// InfoCommands returns a command by id or list of commands by ids
func (a *App) InfoCommands(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("id")
	result, err := a.InfoCommand(param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("error on get info commands ", err)
		return
	}

	if result == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Info("command not found with id ", param)
		return
	}

	jsonCommand, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("error on marshalling info commands ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonCommand)
}

func (a *App) SetupRoutes() http.Handler {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/info", a.InfoAll).Methods("GET")
	a.Router.HandleFunc("/info-by-id", a.InfoCommands).Methods("GET").Queries("id", "{id}")
	a.Router.HandleFunc("/new", a.NewCommand).Methods("POST")
	a.Router.HandleFunc("/remove", a.RemoveCommand).Methods("GET").Queries("id", "{id}")
	return a.Router
}
