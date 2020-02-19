package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/machinable/event-processor/webhook"
	"github.com/sirupsen/logrus"
)

const (
	// QueueHooks is the redis queue for web hooks
	QueueHooks = "hook_queue"
	// QueueHookResults is the redis queue for web hook result messages
	QueueHookResults = "hook_result_queue"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	// create a new redis client
	queue := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PW", ""), // no password set
		DB:       0,                      // use default DB
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
			logger.Error(err)
			return
		}

		// unmarshal event
		event := &webhook.HookEvent{}
		if err := json.Unmarshal([]byte(result[1]), event); err != nil {
			logger.Error(err)
			continue
		}

		// send event
		hookResult := client.PostHook(event)
		b, _ := json.Marshal(hookResult)

		// add hook result to queue
		if err := queue.RPush(QueueHookResults, b).Err(); err != nil {
			logger.Error(err)
		}
	}
}
