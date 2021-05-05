package main

import (
	"fmt"
	"github.com/eladhaim22/sic/services"
	"github.com/eladhaim22/sic/utils/env"
	"github.com/go-co-op/gocron"
	"time"
	log "github.com/sirupsen/logrus"
)

func main(){
	log.Info(fmt.Sprintf("Starting scheduler with cron expression %s", env.CRON))
	s := gocron.NewScheduler(time.UTC)
	s.Cron(env.CRON).Do(services.CleaningTask)
	s.StartBlocking()
}
