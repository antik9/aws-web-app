package main

import (
	"log"
	"os"

	"github.com/antik9/aws-web-app/internal/db"
	"github.com/antik9/aws-web-app/internal/queue"
)

func main() {
	client := queue.NewClient("consumer")
	for {
		ip := client.ReadMessage()

		func(_ip string) {
			f, err := os.OpenFile("/tmp/blacklist", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			if _, err = f.WriteString(_ip + "\n"); err != nil {
				log.Fatal(err)
			}
			db.UpdateBlackList(_ip)
		}(ip)
	}
}
