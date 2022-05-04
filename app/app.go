package main

import (
	"github.com/pieterclaerhout/go-log"
	"redis_channels_playground/app/config"
	"redis_channels_playground/app/metrics"
	"redis_channels_playground/app/repo"
	"redis_channels_playground/app/service"
	"redis_channels_playground/app/web"
)

// start redis
// create a channel
// write to the channel
// read from the channel

func main() {
	log.DebugMode = true
	log.DebugSQLMode = true
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "2006-01-02 15:04:05.000"

	hc := service.NewHealthChecker()

	cfg := config.GetConfig()

	r, err := repo.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Disconnect()

	r2, err := repo.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	search := service.NewService(r, r2)
	hc.RegisterService(r)

	s := web.NewServer(cfg, search, hc)
	if err := hc.Healthy(); err != nil {
		log.Warn("no healthy")
	} else {
		log.Info("service is healthy")
	}

	p := metrics.NewPrometheus(s.Router())
	p.Use(s.Router())

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
