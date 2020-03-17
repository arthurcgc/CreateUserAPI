package data

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type Mdata struct {
	Client     *mongo.Client
	database   string
	collection *mongo.Collection
}

func (db *Mdata) getDbConnectionString() string {
	dbString := os.Getenv("MONGO_ENV")
	return dbString
}

func (db *Mdata) OpenDb() error {
	var err error
	db.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	// Calling Connect does not block for server discovery
	// So to find out if we really connected to the mongo server we call on the Ping method
	err = db.Client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return err
	}

	db.collection = db.getColl("Users")
	return nil
}

func (db *Mdata) getColl(collection string) *mongo.Collection {
	return db.Client.Database("Accounts").Collection(collection)
}

func (db *Mdata) InsertUser(name string, email string) (*User, error) {
	user := &User{
		Name:  name,
		Email: email,
	}
	payload, err := bson.Marshal(user)
	if err != nil {
		return nil, err
	}

	_, err = db.collection.InsertOne(context.TODO(), payload)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Mdata) GetUser(email string) (*User, error) {
	filter := bson.M{"email": email}
	var user User
	err := db.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	//err = bson.Unmarshal(payload, &user)
	//if err != nil {
	//	return nil, err
	//}

	return &user, nil
}

func (db *Mdata) UpdateUser(oldEmail, newEmail, newName string) (*User, error) {
	return nil, nil
}

func (db *Mdata) DeleteUser(email string) (*User, error) {
	return nil, nil
}

func (db *Mdata) GetAll() ([]User, error) {
	return nil, nil
}

func (db *Mdata) CloseDb() error {
	return nil
}
