package main

import (
	"fmt"
	"github.com/benmanns/goworker"
	"github.com/zendesk/slack-poc/jobs"
)

func main() {
	goworker.Register("ProcessSlackMessage", jobs.ProcessSlackMessage)
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}
