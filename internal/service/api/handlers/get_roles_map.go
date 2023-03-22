package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
)

func GetRolesMap(w http.ResponseWriter, r *http.Request) {
	result := newModuleRolesResponse()

	ape.Render(w, result)
}
