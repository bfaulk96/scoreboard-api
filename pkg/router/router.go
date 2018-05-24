package router

import(
	"github.com/gorilla/mux"
	"github.com/bfaulk96/scoreboard-api/pkg/database"
	"github.com/bfaulk96/scoreboard-api/pkg/configurations"
	"github.com/bfaulk96/scoreboard-api/pkg/handlers"
	"net/http"
)

type Router struct {
	*mux.Router
}

func New() (r *Router) {
	return &Router{
		mux.NewRouter().StrictSlash(true),
	}
}

func (r *Router) CreateRoutes(collection *database.DBCollection, config *configurations.Config) () {
	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/games", handlers.GetAllActiveGames(collection)).Methods("GET")
	r.HandleFunc("/games/recent", handlers.GetMostRecentGame(collection)).Methods("GET")
	r.HandleFunc("/games/code/{code}", handlers.GetGameByCode(collection)).Methods("GET")
	r.HandleFunc("/games/name/{name}", handlers.GetGameByName(collection)).Methods("GET")
	r.HandleFunc("/games/new", handlers.StartNewGame(collection)).Methods("POST")
	r.HandleFunc("/games/delete", handlers.DeleteGame(collection)).Methods("DELETE")
	r.HandleFunc("/increment-blue/{increment}", handlers.IncrementBlue(collection, false)).Methods("PUT")
	r.HandleFunc("/increment-blue", handlers.IncrementBlue(collection, true)).Methods("PUT")
	r.HandleFunc("/set-blue/{blueScore}", handlers.SetBlue(collection)).Methods("PUT")
	r.HandleFunc("/increment-red/{increment}", handlers.IncrementRed(collection, false)).Methods("PUT")
	r.HandleFunc("/increment-red", handlers.IncrementRed(collection, true)).Methods("PUT")
	r.HandleFunc("/set-red/{redScore}", handlers.SetRed(collection)).Methods("PUT")

	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundPage)
}