package scripts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func handlerMatchFunc(IdUser int64) bot.MatchFunc {
	return func(update *models.Update) bool {
		if update.Message != nil {
			IdSender := update.Message.From.ID
			if IdSender == IdUser {
				return true
			}
		}
		return false
	}
}

type determineTimeZoneResp struct {
	Timezone string `json:"timezone_id"`
}

func reqForTheTimeZone(chanTimeZone chan<- string, latitude, longitude float64) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.wheretheiss.at/v1/coordinates/%f,%f", latitude, longitude),
		nil,
	)
	if err != nil {
		close(chanTimeZone)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		close(chanTimeZone)
		return
	}
	defer resp.Body.Close()

	var response determineTimeZoneResp
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		close(chanTimeZone)
		return
	}

	chanTimeZone <- response.Timezone
}
