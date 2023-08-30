package models

import "time"

type UserSubscription struct {
	Name    string    `json:"name"`
	TimeOut time.Time `json:"timeOut"`
}

type Subscriber struct {
	UserId int64    `json:"userId"`
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}

type SegmentBody struct {
	Id      int    `json:"id"`
	Name    string `json:"segmentName"`
	Percent int    `json:"percents"`
}

type SubscribeWithTimeout struct {
	UserId      int64  `json:"userId"`
	SegmentName string `json:"segmentName"`
	TimeOut     int    `json:"timeToDie"`
}

type History struct {
	UserId    int64     `json:"userId"`
	Segment   string    `json:"segment"`
	Operation string    `json:"Operation"`
	Time      time.Time `json:"time"`
}
