package main

import (
	"fmt"
	"github.com/eladhaim22/sic/tasks"
	"github.com/eladhaim22/sic/utils/env"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

func main(){
	log.Info(fmt.Sprintf("Starting scheduler with cron expression %s", env.CRON))
	s := gocron.NewScheduler(time.Local)
	s.Cron(env.CRON).Do(tasks.CleaningTask)
	s.StartBlocking()
}
