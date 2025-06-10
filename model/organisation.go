package model

type Organisation struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`

	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}
