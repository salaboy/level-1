package function

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"time"
)

type Questions struct {
	SessionId string
	Question1 string
	Question2 string
	Question3 string
}

type Score struct {
	SessionId string
	Time      time.Time
	Level     string
	LevelScore int
}

type GameTime struct{
	GameTimeId string
	SessionId string
	Level string
	Type string
	Time      time.Time
}

var redisHost = os.Getenv("REDIS_HOST")

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

	points := 0
	var q Questions

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(req.Body).Decode(&q)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if q.Question1 == "42" {
		points++
	} else {
		points--
	}
	if q.Question2 == "43" {
		points++
	} else {
		points--
	}
	if q.Question3 == "44" {
		points++
	} else {
		points--
	}

	var score Score
	score.Level = "level-1"
	score.LevelScore = points
	score.SessionId = q.SessionId
	score.Time = time.Now()
	scoreJson, err := json.Marshal(score)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = client.RPush("score-" + q.SessionId, string(scoreJson)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}


	gt := GameTime{
		GameTimeId: "time-" + score.SessionId,
		SessionId:  score.SessionId,
		Level:      score.Level,
		Type:       "end",
		Time:       score.Time,
	}

	gameTimeJson, err := json.Marshal(gt)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = client.RPush(gt.GameTimeId, string(gameTimeJson)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}


	fmt.Fprintln(res, string(scoreJson))

}
