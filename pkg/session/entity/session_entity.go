package entity

import "time"

type Session struct {
	SessionToken string    `bson:"_id"`
	Username     string    `bson:"uname"`
	Expire       time.Time `bson:"expireAt"`
}
