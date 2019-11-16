package main

import (
	"fmt"
	"os"
	"time"

	"github.com/antik9/aws-web-app/internal/aws"
	"github.com/antik9/aws-web-app/internal/db"
)

func main() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			blackListedIps := db.GetIpsForBlackList()
			for _, ip := range blackListedIps {
				hostname, _ := os.Hostname()
				message := fmt.Sprintf("%s on %s", ip, hostname)
				awsapp.SendMessage(&message)
				db.UpdateBlackList(ip)
			}
		}
	}
}
