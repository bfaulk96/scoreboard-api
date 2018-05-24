package handlers

import (
	"github.com/bfaulk96/scoreboard-api/pkg/database"
	"net/http"
	"github.com/bfaulk96/scoreboard-api/pkg/models"
	"log"
	"github.com/bfaulk96/scoreboard-api/pkg/models/responses"
	"github.com/globalsign/mgo/bson"
	"time"
	"strconv"
	"math/rand"
	"strings"
)

func GetAllActiveGames(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		activeGames, err := collection.GetAllActiveGames()
		if err != nil {
			log.Printf("Mongo error: %v\n", err.Error())
			aw.Respond(ar, &responses.Error{Error: "Unable to fetch active games"}, http.StatusInternalServerError)
			return
		}

		if len(activeGames) == 0 {
			activeGames = make([]models.Scoreboard, 0)
		}

		type Response struct {
			ActiveGames	[]models.Scoreboard	`json:"activeGames"`
		}

		aw.Respond(ar, &Response{ActiveGames: activeGames}, http.StatusOK)
	}
}

func GetMostRecentGame(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		game, err := collection.GetMostRecentGame()
		if err != nil {
			log.Printf("Error finding most recent game: %v\n", err.Error())
			aw.Respond(ar, &responses.Error{Error: "Unable to find most recent game"}, http.StatusNotFound)
			return
		}
		aw.Respond(ar, game, http.StatusOK)
	}
}

func GetGameByName(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		name := ar.GetRouteVariables()["name"]
		game, err := collection.GetActiveGameByName(name)
		if err != nil {
			log.Printf("Error finding game by name: %v\n", err.Error())
			aw.Respond(ar, &responses.Error{Error: "Unable to find game by name"}, http.StatusNotFound)
			return
		}
		aw.Respond(ar, game, http.StatusOK)
	}
}

func GetGameByCode(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		code := ar.GetRouteVariables()["code"]
		game, err := collection.GetActiveGameByCode(code)
		if err != nil {
			log.Printf("Error finding game by code: %v\n", err.Error())
			aw.Respond(ar, &responses.Error{Error: "Unable to find game by code"}, http.StatusNotFound)
			return
		}
		aw.Respond(ar, game, http.StatusOK)
	}
}

func StartNewGame(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		if Authenticate(ar.BasicAuth()) {
			queryParams := ar.GetQueryParameters()
			var err error
			ws := 25
			gn := "Red vs. Blue"
			gc := strings.ToUpper(RandomString(10))
			winningScore, ok := queryParams["winningScore"]
			if ok && len(winningScore) >= 1 {
				ws, err = strconv.Atoi(winningScore[0])
				if err != nil {
					log.Printf("int expected for 'winningScore', actual: %v\n", winningScore)
					aw.Respond(ar, &responses.Error{Error: "winningScore query parameter was not an int"}, http.StatusBadRequest)
					return
				}
			}
			gameName, ok := queryParams["gameName"]
			if ok && len(gameName) >= 1 {
				gn = gameName[0]
			}
			gameCode, ok := queryParams["gameCode"]
			if ok && len(gameCode) >= 1 {
				gc = gameCode[0]
			}

			err = collection.StartNewGame(&models.Scoreboard{
				Id:             bson.NewObjectId(),
				GameStartDate:  time.Now(),
				GameUpdateDate: time.Now(),
				BlueScore:      0,
				RedScore:       0,
				WinningScore:   ws,
				GameName:       gn,
				GameCode:       gc,
				GameFinished:   false,
			})

			if err != nil {
				log.Printf("Error saving new scoreboard to Database: %v\n", err.Error())
				aw.Respond(ar, &responses.Error{Error: "Unable to create new game"}, http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Invalid API Username/Password")
			aw.Respond(ar, &responses.Error{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}
	}
}

func IncrementBlue(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		if Authenticate(ar.BasicAuth()) {
			routeVariables := ar.GetRouteVariables()
			incrementStr := routeVariables["increment"]
			increment, err := strconv.Atoi(incrementStr)
			if err != nil {
				log.Printf("Integer expected for increment value, something else was provided..")
				aw.Respond(ar, &responses.Error{Error: "Integer expected for increment value, something else was provided."}, http.StatusBadRequest)
				return
			}
			game, responseError := FindGame(ar, aw, collection)
			if responseError != nil {
				aw.Respond(ar, responseError, http.StatusNotFound)
				return
			}
			game.BlueScore += increment
			if game.BlueScore >= game.WinningScore {
				game.GameFinished = true
			}
			err = collection.UpdateGame(game.Id.Hex(), game)
		} else {
			log.Printf("Invalid API Username/Password")
			aw.Respond(ar, &responses.Error{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}
	}
}

func IncrementRed(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		if Authenticate(ar.BasicAuth()) {
			routeVariables := ar.GetRouteVariables()
			incrementStr := routeVariables["increment"]
			increment, err := strconv.Atoi(incrementStr)
			if err != nil {
				log.Printf("Integer expected for increment value, something else was provided.")
				aw.Respond(ar, &responses.Error{Error: "Integer expected for increment value, something else was provided."}, http.StatusBadRequest)
				return
			}
			game, responseError := FindGame(ar, aw, collection)
			if responseError != nil {
				aw.Respond(ar, responseError, http.StatusNotFound)
				return
			}
			game.RedScore += increment
			if game.RedScore >= game.WinningScore {
				game.GameFinished = true
			}
			err = collection.UpdateGame(game.Id.Hex(), game)
		} else {
			log.Printf("Invalid API Username/Password")
			aw.Respond(ar, &responses.Error{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}
	}
}

func SetBlue(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		if Authenticate(ar.BasicAuth()) {
			routeVariables := ar.GetRouteVariables()
			blueScoreStr := routeVariables["blueScore"]
			blueScore, err := strconv.Atoi(blueScoreStr)
			if err != nil {
				log.Printf("Integer expected for blue score, something else was provided..")
				aw.Respond(ar, &responses.Error{Error: "Integer expected for blue score, something else was provided."}, http.StatusBadRequest)
				return
			}
			game, responseError := FindGame(ar, aw, collection)
			if responseError != nil {
				aw.Respond(ar, responseError, http.StatusNotFound)
				return
			}
			game.BlueScore = blueScore
			if blueScore >= game.WinningScore {
				game.GameFinished = true
			}
			err = collection.UpdateGame(game.Id.Hex(), game)
			if err != nil {
				log.Printf("Error Updating scoreboard: %v\n", err.Error())
				aw.Respond(ar, &responses.Error{Error: "Unable to update scoreboard"}, http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Invalid API Username/Password")
			aw.Respond(ar, &responses.Error{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}
	}
}

func SetRed(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		if Authenticate(ar.BasicAuth()) {
			routeVariables := ar.GetRouteVariables()
			redScoreStr := routeVariables["redScore"]
			redScore, err := strconv.Atoi(redScoreStr)
			if err != nil {
				log.Printf("Integer expected for red score, something else was provided..")
				aw.Respond(ar, &responses.Error{Error: "Integer expected for red score, something else was provided."}, http.StatusBadRequest)
				return
			}
			game, responseError := FindGame(ar, aw, collection)
			if responseError != nil {
				aw.Respond(ar, responseError, http.StatusNotFound)
				return
			}
			game.RedScore = redScore
			if redScore >= game.WinningScore {
				game.GameFinished = true
			}
			err = collection.UpdateGame(game.Id.Hex(), game)
			if err != nil {
				log.Printf("Error Updating scoreboard: %v\n", err.Error())
				aw.Respond(ar, &responses.Error{Error: "Unable to update scoreboard"}, http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Invalid API Username/Password")
			aw.Respond(ar, &responses.Error{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}
	}
}

func DeleteGame(collection *database.DBCollection) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) () {
		ar, aw := models.NewAPIRequest(r), models.NewAPIResponseWriter(w)
		if Authenticate(ar.BasicAuth()) {
			game, responseError := FindGame(ar, aw, collection)
			if responseError != nil {
				aw.Respond(ar, responseError, http.StatusNotFound)
				return
			}

			err := collection.DeleteGame(game.Id.Hex())
			if err != nil {
				log.Printf("Error Deleting scoreboard: %v\n", err.Error())
				aw.Respond(ar, &responses.Error{Error: "Unable to delete scoreboard"}, http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Invalid API Username/Password")
			aw.Respond(ar, &responses.Error{Error: "Unauthorized"}, http.StatusUnauthorized)
			return
		}
	}
}


func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func Authenticate(user string, pass string, ok bool) (bool) {
	ExpectedUser := "scoreboardUser"
	ExpectedPass := "BgF4QnYhytgpMFPYfTQ2cWprxhUvaVwz4XYc6CzkCv5Ss2Ug"
	return ok && user == ExpectedUser && pass == ExpectedPass
}

func FindGame(ar *models.APIRequest, aw *models.APIResponseWriter, collection *database.DBCollection) (game *models.Scoreboard, responseError *responses.Error){
	var ok bool
	queryParams := ar.GetQueryParameters()
	var gameName []string
	var gameCode []string
	gameCode, ok = queryParams["gameCode"]
	var err error
	if ok && len(gameCode) >= 1 {
		game, err = collection.GetActiveGameByCode(gameCode[0])
		if err != nil {
			log.Printf("Error finding game by code: %v\n", err.Error())
			return nil, &responses.Error{Error: "Unable to find game using code"}
		}
	}  else if gameCode, ok = queryParams["code"]; ok && len(gameCode) >= 1 {
		game, err = collection.GetActiveGameByCode(gameCode[0])
		if err != nil {
			log.Printf("Error finding game by code: %v\n", err.Error())
			return nil, &responses.Error{Error: "Unable to find game using code"}
		}
	} else if gameName, ok = queryParams["gameName"]; ok && len(gameName) >= 1 {
		game, err = collection.GetActiveGameByName(gameName[0])
		if err != nil {
			log.Printf("Error finding game by name: %v\n", err.Error())
			return nil, &responses.Error{Error: "Unable to find game using name"}
		}
	}  else if gameName, ok = queryParams["name"]; ok && len(gameName) >= 1 {
		game, err = collection.GetActiveGameByName(gameName[0])
		if err != nil {
			log.Printf("Error finding game by name: %v\n", err.Error())
			return nil, &responses.Error{Error: "Unable to find game using name"}
		}
	} else {
		game, err = collection.GetMostRecentGame()
		if err != nil {
			log.Printf("Error finding most recent game: %v\n", err.Error())
			return nil, &responses.Error{Error: "Unable to find most recent game"}
		}
	}
	return game, nil
}
