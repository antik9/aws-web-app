package main

import (
	"time"

	"github.com/antik9/aws-web-app/internal/db"
	"github.com/antik9/aws-web-app/internal/queue"
)

func main() {
	ticker := time.NewTicker(time.Second * 5)
	client := queue.NewClient("producer")
	for {
		select {
		case <-ticker.C:
			blackListedIps := db.GetIpsForBlackList()
			for _, ip := range blackListedIps {
				client.SendMessage(ip)
			}
		}
	}
}
