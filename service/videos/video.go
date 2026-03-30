package videos

import "go.mongodb.org/mongo-driver/v2/bson"

type video struct {
	ID   bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Path string        `json:"path" bson:"path"`
}
