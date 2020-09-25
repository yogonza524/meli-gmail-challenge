package main

import (
	"log"
	"meli/domain/gmail"
	"meli/domain/sql"
)

func main() {
	findEmails()
}

func findEmails() {
	srv, err := gservice.CreateService()
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	db := dao.Connect()
	gservice.FindMessages("{from:-devops} AND {subject:devops} OR {devops}", srv, db)
}
