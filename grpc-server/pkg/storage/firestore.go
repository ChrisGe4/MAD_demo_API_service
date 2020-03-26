package storage

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"

	"github.com/chrisge4/MAD_demo_API_service/grpc-server/pkg/pb"
)

type Firestore struct {
	collection string
	client     *firestore.Client
}

func NewFirestore(ctx context.Context, credentialsFile, projectId, collection string) (*Firestore, error) {
	log.Printf("Initializing Firestore with collection %q in projet %q", collection, projectId)
	var client *firestore.Client
	var err error
	if credentialsFile == "" {
		client, err = firestore.NewClient(ctx, projectId)
	} else {
		client, err = firestore.NewClient(ctx, projectId, option.WithCredentialsFile(credentialsFile))

	}

	if err != nil {
		return nil, err
	}

	return &Firestore{client: client, collection: collection}, nil
}

func (fs *Firestore) Add(ctx context.Context, t *pb.Task) error {

	_, err := fs.client.Collection(fs.collection).Doc(t.Id).Set(ctx, t)
	if err != nil {
		return err
	}
	return nil
}
func (fs *Firestore) List(ctx context.Context, category *pb.Category) ([]*pb.Task, error) {
	ti := fs.client.Collection(fs.collection).Where("Category.Name", "==", category.Name).Documents(ctx)
	var tasks []*pb.Task
	for {
		doc, err := ti.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		jb, err := json.Marshal(doc.Data())
		if err != nil {
			return nil, err
		}

		var task pb.Task
		if err := json.Unmarshal(jb, &task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (fs *Firestore) Close() error {
	return fs.client.Close()
}
