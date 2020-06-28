package main

import (
	"flag"
	"log"

	"github.com/robfig/cron/v3"
)

var (
	flagJobsDir = flag.String("jobs-dir", ".", "the jobs directory")
)

var (
	jobs    map[string]*Job
	cronhub *cron.Cron
)

func init() {
	flag.Parse()

	var err error

	cronhub = cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	jobs, err = ParseJobs(*flagJobsDir)
	if err != nil {
		log.Fatal(err.Error())
	}

}
