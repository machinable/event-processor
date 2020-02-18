package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/machinable/event-processor/webhook"
	"github.com/sirupsen/logrus"
)

func main() {
	// create a new redis client
	queue := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// ping the redis server
	pong, err := queue.Ping().Result()

	// initialize logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	logger.Info("starting event processor")

	// initialize client
	client := webhook.NewClient(logger)

	// fail quickly if ping fails
	if err != nil {
		fmt.Println(pong, err)
		return
	}

	logger.Info("waiting for events...")
	for {
		// endlessly read from queue
		result, err := queue.BLPop(0, "hook_queue").Result()

		// exit on a read error
		if err != nil {
			log.Println(err)
			return
		}

		// unmarshal event
		event := &webhook.HookEvent{}
		if err := json.Unmarshal([]byte(result[1]), event); err != nil {
			logger.Error(err)
			continue
		}

		// send event
		client.PostHook(event)
	}
}
