package model

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"time"
)

type RssService struct {
	ID          int64
	Name        string `gorm:"unique_index"`
	Url         string `gorm:"unique_index"`
	Title       string
	Description string
	Link        string
	Language    string
	NumFollwer  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func AllServices() []RssService {
	var services []RssService
	if err := DB.Find(&services).Error; err != nil {
		log.Print(err)
		return []RssService{}
	}
	return services
}

func GetServiceByName(name string) (RssService, error) {
	var service RssService
	if err := DB.First(&service, "name = ?", name).Error; err != nil {
		log.Print(err)
		return RssService{}, err
	}
	return service, nil
}

func CreateService(newServ RssService) error {
	if err := DB.Create(&newServ).Error; err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func DeleteServiceByName(name string) error {
	if err := DB.Where("service_name = ?", name).Delete(&Subs{}).Error; err != nil {
		return err
	}
	if err := DB.Where("name = ?", name).Delete(&RssService{}).Error; err != nil {
		log.Print(err)
		return err
	}
	return nil
}

type RssItem struct {
	ID          int64 `json:"-"`
	ServiceName string
	Title       string
	Link        string `gorm:"unique_index"`
	Published   string
	CreatedAt   time.Time `json:"-"`
}

func CreateRssItems(feed *gofeed.Feed, service RssService, addRedis bool) int {
	isFirstAdd := true
	numAddItems := 0
	for _, i := range feed.Items {
		newItem := RssItem{
			ServiceName: service.Name,
			Title:       i.Title,
			Link:        i.Link,
			Published:   i.Published,
		}
		if err := DB.Create(&newItem).Error; err != nil {
			break
		}
		msg, _ := json.Marshal(newItem)
		Logger.Write("RssItem", newItem)
		numAddItems += 1
		if !addRedis {
			continue
		}
		Redis.LPush(fmt.Sprintf("service-%s", service.Name), msg)
		if isFirstAdd {
			Redis.SAdd("newItem", service.Name)
			isFirstAdd = false
		}
	}
	return numAddItems
}
