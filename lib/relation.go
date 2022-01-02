package lib

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const coll = "relations"

type Relation struct {
	ProjectID string `json:"projectID"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Kind      string `json:"kind,omitempty"`
	Location  string `json:"location,omitempty"`
}

func (r Relation) Upload() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDBTransactionTimeout)
	defer cancel()

	result, err := db.Collection(coll).InsertOne(ctx, &r)
	if err != nil {
		return errors.WithStack(err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("failed to convert InsertedID to string")
	}

	fmt.Printf("ID=%s SRC=%s TAR=%s\n", id.Hex(), r.Source, r.Target)

	return nil
}

func Parse(input string) (Relation, error) {
	var relation Relation
	if err := json.Unmarshal([]byte(input), &relation); err != nil {
		return relation, err
	}

	return relation, nil
}
