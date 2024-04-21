package utils

import (
	"fmt"

	"github.com/Deepjyoti-Sarmah/GolangSearchEngine/search"
	"github.com/robfig/cron/v3"
) 

func StartCornJobs() {
  c := cron.New()
  c.AddFunc("0 * * * *", search.RunEngine) //Run every hour
  c.Start()
  cronCount := len(c.Entries())
  fmt.Printf("setup %d corn jobs \n", cronCount)
}

func runEngine()  {
  fmt.Println("Starting engine")
}
