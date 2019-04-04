package handlers

import (
	"html/template"
	"log"

	"github.com/mchmarny/tevents/pkg/utils"
)

const (
	defaultPublisherToken = "notset"
)

var (
	// Templates for handlers
	templates           *template.Template
	knownPublisherToken string
)

// InitHandlers initializes OAuth package
func InitHandlers() {

	// Templates
	tmpls, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error while parsing templates: %v", err)
	}
	templates = tmpls

	// know publisher
	knownPublisherToken = utils.MustGetEnv("KNOWN_PUBLISHER_TOKEN", defaultPublisherToken)

	initWS()

}
