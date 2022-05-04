package repo

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/pieterclaerhout/go-log"
	"math/rand"
	"redis_channels_playground/app/config"
	"redis_channels_playground/app/dto"
	"strconv"
	"strings"
	"time"
)

const consumer = "consumer"

type Redis struct {
	client     *redis.Client
	streamName string
	id         string
	group      string
}

func NewRedis(cfg config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHostPort,
		Password: "",
		DB:       0,
	})

	r := &Redis{
		client:     client,
		streamName: cfg.StreamName,
		id:         strconv.Itoa(rand.Int()),
		group:      strconv.Itoa(rand.Int()),
	}

	err := r.RegisterGroup()
	if strings.Contains(err.Error(), "BUSYGROUP") {
		err = nil
	}
	return r, err
}

func (r *Redis) HealthCheck() error {
	c := r.client.Ping()
	return c.Err()
}

func (r *Redis) Disconnect() {
	_ = r.client.Close()
}

func (r *Redis) Push(m *dto.Message) error {
	d, err := json.Marshal(m)
	if err != nil {
		return err
	}
	cmd := r.client.XAdd(&redis.XAddArgs{
		Stream: r.streamName,
		Values: map[string]interface{}{m.Name: d},
	})
	return cmd.Err()
}

func (r *Redis) RegisterGroup() error {
	return r.client.XGroupCreate(r.streamName, r.group, "$").Err()
}

func (r *Redis) Pull() error {
	resp, err := r.client.XReadGroup(&redis.XReadGroupArgs{
		Group:    r.group,
		Consumer: r.id,
		Streams:  []string{r.streamName, ">"},
		Block:    time.Minute,
	}).Result()
	if err != nil {
		return err
	}

	for streamName, messages := range resp {
		for _, m := range messages.Messages {
			log.Debugf("%v %v %v %v", r.id, streamName, m.ID, m.Values)
		}
	}
	return nil
}
