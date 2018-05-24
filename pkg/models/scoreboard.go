package models

import (
	"time"
	"github.com/globalsign/mgo/bson"
)

type Scoreboard struct {
	Id				bson.ObjectId	`json:"_id" bson:"_id,omitempty"`
	GameStartDate	time.Time 		`json:"gameStartDate" bson:"gameStartDate"`
	GameUpdateDate	time.Time 		`json:"gameUpdateDate" bson:"gameUpdateDate"`
	BlueScore		int				`json:"blueScore" bson:"blueScore"`
	RedScore		int				`json:"redScore" bson:"redScore"`
	WinningScore	int				`json:"winningScore" bson:"winningScore"`
	GameName		string			`json:"gameName" bson:"gameName"`
	GameCode		string			`json:"gameCode" bson:"gameCode"`
	GameFinished	bool			`json:"gameFinished" bson:"gameFinished"`
}