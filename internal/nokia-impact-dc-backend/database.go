package nokia_impact_dc_backend

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"net/http"
)

type Database struct {
	client           *firestore.Client
	displayNameCache map[string]string
}

func (d *Database) log(request *http.Request, s string) {
	// TODO
	//d.client.Collection("logs").Doc(request.Tim)
}

func (d *Database) saveData(r *http.Request, record *AbstractDataRecord) error {
	client := GetFirestoreClient(r)
	docId := fmt.Sprintf("%d", record.Timestamp)
	if record.SerialNumber != "" {
		ctx := context.Background()
		recordTyp := struct {
			RecordType string `firestore:"recordType"`
		}{ record.UpdateType}
		docRef := client.Collection("clients").Doc(record.SerialNumber).Collection("history").Doc(docId)
		if _, err := docRef.Set(ctx, recordTyp); err != nil {
			return err
		}
		if _, err := docRef.Set(ctx, record.StringValues, firestore.MergeAll); err != nil {
			return err
		}
		if _, err := docRef.Set(ctx, record.NumberValues, firestore.MergeAll); err != nil {
			return err
		}
		if _, err := docRef.Set(ctx, record.BooleanValues, firestore.MergeAll); err != nil {
			return err
		}
		if _, err := docRef.Set(ctx, record.ArrayValues, firestore.MergeAll); err != nil {
			return err
		}
	}
	return nil
}

func NewDatabase(cli *firestore.Client) *Database {
	return &Database{
		client: cli,
		displayNameCache: make(map[string]string),
	}
}


