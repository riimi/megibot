package main

import (
	_ "fmt"
	"megibot/model"
	"log"
)

func main() {
	model.InitDB("root:1233@tcp(35.200.67.123:3306)/bot?parseTime=true&&charset=utf8mb4,utf8")
	defer model.DB.Close()
	model.InitLogger("localhost:5170")
	defer model.Logger.Close()
	type Result struct {
		ServiceName string
		Count       int
	}
	var res []Result
	if err := model.DB.Model(model.RssItem{}).Select("service_name, count(*) as count").Group("service_name").Having("count(*) > ?", 100).Scan(&res).Error; err != nil {
		log.Fatal(err)
	}
	//fmt.Println(res)
	model.Logger.Write("test", res[0])
	/*
		for _, r := range res {
			if err := model.DB.Exec("DELETE FROM rss_items WHERE service_name=? AND id < (SELECT id FROM (SELECT id FROM rss_items WHERE service_name=? ORDER BY id DESC LIMIT 1 OFFSET 100) foo)", r.ServiceName, r.ServiceName).Error; err != nil {
				log.Fatal(err)
			}
		}
	*/
}
