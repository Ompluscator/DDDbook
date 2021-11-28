package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/PacktPublishing/Domain-Driven-Design-with-Go/chapter5/internal/user"
)

func main() {
	// create all third-party dependencies
	validate := validator.New()
	engine := gin.Default()

	// creates new user module
	module := &user.Module{}
	module.Configure("./database.sqlite", engine, validate)

	// run server
	log.Fatalln(engine.Run())
}
