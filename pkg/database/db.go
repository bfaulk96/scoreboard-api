package database

import (
	"github.com/globalsign/mgo"
	"github.com/bfaulk96/scoreboard-api/pkg/configurations"
	"log"
	"strings"
	"encoding/base64"
	"github.com/bfaulk96/scoreboard-api/pkg/models"
	"github.com/globalsign/mgo/bson"
	"errors"
)

type DBCollection struct {
	Collection *mgo.Collection
}

func InitializeMongoDatabase(config *configurations.Config) (collection *DBCollection) {
	var connectionString string
	if config.DbURL != "" {
		dbURL := config.DbURL
		if strings.Contains(strings.ToLower(dbURL), "mongodb") {
			connectionString = dbURL
		} else {
			url, err := base64.StdEncoding.DecodeString(dbURL)
			if err != nil {
				log.Fatal("Error base64 decoding connection string")
			}
			connectionString = string(url)
		}
	} else {
		log.Fatal("No DB Connection string provided in Config file.")
	}

	session, err := mgo.Dial(connectionString)
	if err != nil {
		//println(connectionString)
		log.Fatalf("Error connecting to database:\n%v\n", err)
	}
	session.SetMode(mgo.Monotonic, true)
	mgoCollection := session.DB(config.DbName).C(config.CollectionName)
	return &DBCollection{Collection: mgoCollection}
}

func (collection *DBCollection) GetAllActiveGames() (activeGames []models.Scoreboard, err error) {
	err = collection.Collection.Find(bson.M{"gameFinished": false}).All(&activeGames)
	return activeGames, err
}

func (collection *DBCollection) StartNewGame(newGame *models.Scoreboard) (err error) {
	return collection.Collection.Insert(&newGame)
}

func (collection *DBCollection) GetActiveGameByName(gameName string) (activeGame *models.Scoreboard, err error) {
	err = collection.Collection.Find(bson.M{"gameName": gameName}).Sort("-gameUpdateDate").One(&activeGame)
	return activeGame, err
}

func (collection *DBCollection) GetActiveGameByCode(gameCode string) (activeGame *models.Scoreboard, err error) {
	err = collection.Collection.Find(bson.M{"gameCode": gameCode}).Sort("-gameUpdateDate").One(&activeGame)
	return activeGame, err
}

func (collection *DBCollection) GetMostRecentGame() (activeGame *models.Scoreboard, err error) {
	err = collection.Collection.Find(bson.M{"gameFinished": false}).Sort("-gameUpdateDate").One(&activeGame)
	return activeGame, err
}

func (collection *DBCollection) UpdateGame(id string, scoreboard *models.Scoreboard) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}
	return collection.Collection.UpdateId(bson.ObjectIdHex(id), scoreboard)
}

func (collection *DBCollection) DeleteGame(id string) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}
	return collection.Collection.RemoveId(bson.ObjectIdHex(id))
}
