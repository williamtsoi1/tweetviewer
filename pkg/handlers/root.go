package handlers

import (
	"log"
	"net/http"

	"github.com/mchmarny/tevents/pkg/utils"
)

// RootHandler handles view page
func RootHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Serving index... ")

	proto := r.Header.Get("x-forwarded-proto")
	if proto == "" {
		proto = "http"
	}

	data := make(map[string]interface{})

	data["host"] = r.Host
	data["proto"] = proto
	data["version"] = utils.MustGetEnv("RELEASE", "v0-not-set")

	log.Printf("data: %v", data)

	// anonymous
	if err := templates.ExecuteTemplate(w, "index", data); err != nil {
		log.Printf("Error in home template: %s", err)
	}

}
