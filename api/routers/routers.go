package routers

import (
	"kiln-projects/api/handlers"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Implémentation de l'interface http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func InitRoutes(mux *http.ServeMux) {
	// Possibilité d'ajouter des routes supplémentaires sans dénaturer le code
	mux.HandleFunc("GET /xtz/delegations", handlers.GetDelegations)
}
