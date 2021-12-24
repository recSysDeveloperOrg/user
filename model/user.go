package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sync"
)

type UserDao struct {
}

type UserWithID struct {
	User
	ID string `bson:"_id"`
}

type User struct {
	Gender           string `bson:"gender"`
	LastRefreshToken string `bson:"last_refresh_token"`
	Password         string `bson:"password"`
	Username         string `bson:"user_name"`
}

var userDao *UserDao
var userDaoOnce sync.Once

var ErrDuplicateName = errors.New("duplicate username")

const errDuplicateCode = 11000

func NewUserDao() *UserDao {
	userDaoOnce.Do(func() {
		userDao = &UserDao{}
	})

	return userDao
}

func (*UserDao) QueryByID(ctx context.Context, id string) (*UserWithID, error) {
	userObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user UserWithID
	if err := GetClient().Collection(CollectionUser).FindOne(ctx, bson.D{{"_id", userObjectID}}).
		Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (*UserDao) QueryByUsernamePassword(ctx context.Context, username, password string) (*UserWithID, error) {
	var user UserWithID
	if err := GetClient().Collection(CollectionUser).FindOne(ctx, bson.D{{"user_name", username},
		{"password", password}}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("no user for username:%s password:%s", username, password)
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (*UserDao) UpdateRefreshToken(ctx context.Context, userID, refreshToken string) (int64, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}
	ur, err := GetClient().Collection(CollectionUser).UpdateOne(ctx, bson.D{{"user_id", userObjectID}},
		bson.D{{"$set", bson.D{{"last_refresh_token", refreshToken}}}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, err
		}
	}

	return ur.ModifiedCount, nil
}

func (*UserDao) InsertNewUser(ctx context.Context, user *User) error {
	_, err := GetClient().Collection(CollectionUser).InsertOne(ctx, user)
	if err != nil {
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			if len(writeErr.WriteErrors) > 0 {
				if writeErr.WriteErrors[0].Code == errDuplicateCode {
					return ErrDuplicateName
				}
			}
		}
		return err
	}

	return nil
}
