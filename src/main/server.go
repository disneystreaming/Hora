// package main is the entrypoint into the Hora application
package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/disneystreaming/Hora/src/main/docs"
	"github.com/disneystreaming/Hora/src/routes"
)

const (
	// the default serving port
	defaultPort string = "8080"
	// the default serving host
	defaultHost string = "localhost"

	// the name of the environment variable holding the serving port
	portEnvVarName string = "PORT"
)

// @title Hora
// @description Hora is a payload validation service

// @contact.name Your Team
// @contact.url your-team.your-team.com
// @contact.email your-team@your-team.com

// HoraVersion is passed in on build
var HoraVersion string

// configureSwagger accepts the application version, host and port and initializes the documentation metadata
func configureSwagger(ver string, host string, port string) {
	docs.SwaggerInfo.Version = ver
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, port)
}

// main is the entrypoint into the application and serves the sample page
func main() {
	// init.
	s := gin.Default()
	routes.AddRoutes(s)

	// configure
	port, exists := os.LookupEnv(portEnvVarName)
	if !exists {
		port = defaultPort
	}
	configureSwagger(HoraVersion, defaultHost, port)

	// run
	err := s.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
