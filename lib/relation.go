package lib

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

const coll = "relations"

type Relation struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Kind     string `json:"kind,omitempty"`
	Location string `json:"location,omitempty"`
}

func (r Relation) Upload() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDBTransactionTimeout)
	defer cancel()

	_, err := db.Collection(coll).InsertOne(ctx, &r)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Parse(input string) (Relation, error) {
	var relation Relation
	if err := json.Unmarshal([]byte(input), &relation); err != nil {
		return relation, err
	}

	return relation, nil
}
