package backend

import (
	"cloud.google.com/go/firestore"
)

type Database struct {
	client           *firestore.Client
	displayNameCache map[string]string
}

func NewDatabase(cli *firestore.Client) *Database {
	return &Database{
		client: cli,
		displayNameCache: make(map[string]string),
	}
}


