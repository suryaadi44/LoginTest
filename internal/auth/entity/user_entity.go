package entity

type User struct {
	Username string `bson:"_id"`
	Password string `bson:"pass"`
}
