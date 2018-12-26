package model

import (
	"log"
	"time"
)

type Subs struct {
	ID          int64
	UserId      string
	ServiceName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func AllSubs() []Subs {
	var subs []Subs
	if err := DB.Find(&subs).Error; err != nil {
		log.Print(err)
		return []Subs{}
	}
	return subs
}

func AllSubsByUser(id string) ([]Subs, error) {
	var subs []Subs
	if err := DB.Where("user_id = ?", id).Find(&subs).Error; err != nil {
		return []Subs{}, err
	}
	return subs, nil
}

func Subscribe(userId, serviceName string) error {
	if err := DB.Create(&Subs{UserId: userId, ServiceName: serviceName}).Error; err != nil {
		return err
	}
	return nil
}

func Unsubscribe(userId, serviceName string) error {
	err := DB.Where("service_name=? AND user_id=?", serviceName, userId).Delete(&Subs{}).Error
	if err != nil {
		return err
	}
	return nil
}

func Subscriber(serviceName string) []Subs {
	var subs []Subs
	if err := DB.Where("service_name = ?", serviceName).Find(&subs).Error; err != nil {
		log.Print(err)
		return []Subs{}
	}
	return subs
}
