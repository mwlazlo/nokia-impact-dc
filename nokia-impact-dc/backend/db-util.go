package backend

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
)

func safeExtractString(doc *firestore.DocumentSnapshot, key string, dflt string) (string) {
	if value, ok := doc.Data()[key]; !ok {
		log.Printf("Key '%s' not found in document %s\n", key, doc.Ref.Path)
		return dflt
	} else {
		if rv, ok := value.(string); !ok {
			log.Printf("Key '%s'[%s] not a string in document %s\n", key, value, doc.Ref.Path)
			return dflt
		} else {
			return rv
		}
	}
}

func safeExtractBool(doc *firestore.DocumentSnapshot, key string, dflt bool) (bool) {
	if value, ok := doc.Data()[key]; !ok {
		log.Printf("Key '%s' not found in document %s\n", key, doc.Ref.Path)
		return dflt
	} else {
		if rv, ok := value.(bool); !ok {
			log.Printf("Key '%s'[%s] not a boolean in document %s\n", key, value, doc.Ref.Path)
			return dflt
		} else {
			return rv
		}
	}
}

func (d *Database) deleteTree(ref *firestore.DocumentRef) error {
	batch := d.client.Batch()
	if err := d.deleteTreeRecursive(ref, batch); err != nil {
		return err
	}
	if writeResult, err := batch.Commit(context.Background()); err != nil {
		for _, e := range writeResult {
			log.Printf("deleteTree [%s]: %s", ref.Path, e)
		}
		return err
	}
	return nil
}

func (d *Database) deleteTreeRecursive(ref *firestore.DocumentRef, batch *firestore.WriteBatch) error {
	iter := ref.Collections(context.Background())
	for {
		c, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		docIter := c.Documents(context.Background())
		for {
			doc, err := docIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			if err := d.deleteTreeRecursive(doc.Ref, batch); err != nil {
				return err
			}
		}
	}
	batch.Delete(ref)
	return nil
}
