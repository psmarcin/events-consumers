package main

import (
	"log"

	_ "events-consumers/admin/pkg/config"
	server "events-consumers/admin/pkg/http"
	"events-consumers/admin/pkg/jobs"
)

func main() {
	job, err := jobs.New()
	if err != nil {
		log.Fatal(err)
	}

	deps := server.Dependencies{
		Job: job,
	}

	server.Start(deps)
}
