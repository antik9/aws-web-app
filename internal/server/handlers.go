package server

import (
	"fmt"
	"net/http"

	"github.com/antik9/aws-web-app/internal/db"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-FORWARDED-FOR")
	fmt.Fprintf(w, "Hello! Your ip is %s", ip)
	db.SaveIpToDatabase(ip)
}
