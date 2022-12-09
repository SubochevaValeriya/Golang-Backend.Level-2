package framework_level

import (
	"myPackage/internal/application_level"
	"net/http"
)

type Handler struct {
	app *application_level.App
}

func (h *Handler) AddNewUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ChangeEnvironment(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SearchUserByUserOrEnv(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SearchEnvByUserOrEnv(w http.ResponseWriter, r *http.Request) {

}
