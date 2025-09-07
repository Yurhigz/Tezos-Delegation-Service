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

// Implémentation de l'interface http.Handler pour utiliser le router directement en tant que paramètre du serveur http
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) InitRoutes() {
	// Possibilité d'ajouter des routes supplémentaires sans dénaturer le code
	r.mux.HandleFunc("GET /xtz/delegations", handlers.GetDelegations)
}
