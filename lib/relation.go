package lib

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const coll = "relations"

type RelationByTarget struct {
	Language       string `json:"language"`
	TargetModule   string `json:"targetModule"`
	TargetFunc     string `json:"targetFunc,omitempty"`
	SourceModule   string `json:"sourceModule"`
	SourceLocation string `json:"sourceLocation,omitempty"`
}

// db.getCollection("relations").createIndex({ projectID: 1, language: 1, targetModule: 1 }, { unique: true })
type Relation struct {
	ProjectID    string `json:"projectID"`
	Language     string `json:"language"`
	TargetModule string `json:"targetModule"`

	Calls []Call `json:"calls"`
}

type Call struct {
	SourceModule   string `json:"sourceModule"`
	SourceLocation string `json:"sourceLocation,omitempty"`
	TargetFunc     string `json:"targetFunc,omitempty"`
}

func (r RelationByTarget) Upload(projectID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDBTransactionTimeout)
	defer cancel()

	opts := options.Update()
	opts.SetUpsert(true)
	result, err := db.Collection(coll).UpdateOne(ctx, bson.M{
		"projectID":    projectID,
		"targetModule": r.TargetModule,
	}, bson.M{
		"$set": bson.M{
			"projectID":    projectID,
			"language":     r.Language,
			"targetModule": r.TargetModule,
		},
		"$addToSet": bson.M{
			"calls": Call{
				SourceModule:   r.SourceModule,
				SourceLocation: r.SourceLocation,
				TargetFunc:     r.TargetFunc,
			},
		},
	}, opts)
	// result, err := db.Collection(coll).InsertOne(ctx, &r)
	if err != nil {
		return errors.WithStack(err)
	}
	id := "Updated"
	if result.UpsertedCount > 0 {
		if objectID, ok := result.UpsertedID.(primitive.ObjectID); ok {
			id = fmt.Sprintf("Added:%s", objectID.Hex())
		} else {
			id = "Error"
		}
	}

	fmt.Printf("ID=%s SRC=%s TAR=%s\n", id, r.SourceModule, r.TargetModule)

	return nil
}

func Parse(input string) (RelationByTarget, error) {
	var relation RelationByTarget
	if err := json.Unmarshal([]byte(input), &relation); err != nil {
		return relation, err
	}

	return relation, nil
}
