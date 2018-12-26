package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"kitchensink/model"
	_ "kitchensink/template"
	"log"
	"os/exec"
	"strings"
	_ "time"
)

type CmdManager struct {
	CmdCh chan Cmd
}

type Cmd struct {
	Value      string
	Done       chan string
	Message    string
	ReplyToken string
	Source     *linebot.EventSource
}

var rootCmd = &cobra.Command{
	Use: "megibot",
}

var (
	Manager    CmdManager
	replyToken string
	source     *linebot.EventSource
	bot        *linebot.Client
	MegiId     string = "U80c288156ed29a6cfa61e8325df0e271c"
)

func (m *CmdManager) Work(client *linebot.Client) {
	bot = client
	model.InitRedis("localhost:6379")
	defer model.Redis.Close()
	model.InitDB("root:1233@tcp(35.200.67.123:3306)/bot?parseTime=true&&charset=utf8mb4,utf8")
	defer model.DB.Close()
	model.DB.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(
		&model.RssService{}, &model.RssItem{}, &model.User{}, &model.Subs{})
	model.DB.Model(&model.Subs{}).AddUniqueIndex("unique_userid_serviceid", "service_name", "user_id")
	model.InitLogger("localhost:5170")
	defer model.Logger.Close()

	c := cron.New()
	c.AddFunc("0 1-59/10 * * * *", func() {
		publishRssItems()
	})
	c.AddFunc("0 */10 * * * *", func() {
		refreshRssItems()
	})
	c.AddFunc("0 0 10,17 * * 4", func() {
		NotifyWeeklyReport()
	})
	c.AddFunc("0 3 0 * * *", func() {
		cleanRssItems()
	})
	c.Start()

	m.CmdCh = make(chan Cmd, 100)

	for {
		select {
		case in := <-m.CmdCh:
			replyToken = in.ReplyToken
			source = in.Source

			args := strings.Fields(in.Message)
			_, _, err := rootCmd.Find(args)
			if args[0] != "help" && err != nil {
				continue
			}
			rootCmd.SetArgs(args)
			var buf bytes.Buffer
			rootCmd.SetOutput(&buf)
			if err := rootCmd.Execute(); err != nil {
				log.Print(err)
			}
			in.Done <- buf.String()
		}
	}
}

func (m *CmdManager) SendCmd(c Cmd) <-chan string {
	m.CmdCh <- c
	return c.Done
}

func refreshRssItems() {
	var Services []model.RssService
	if err := model.DB.Find(&Services).Error; err != nil {
		log.Print(err)
		return
	}

	for _, s := range Services {
		fp := gofeed.NewParser()
		out, err := exec.Command("curl", "-c", "cookie", "-XGET", "-L", s.Url).Output()
		if err != nil {
			log.Print("curl:" + err.Error())
			continue
		}
		feed, err := fp.ParseString(string(out))
		if err != nil {
			log.Print("parse:" + err.Error())
			continue
		}
		model.CreateRssItems(feed, s, true)
	}
}

func publishRssItems() {
	services, err := model.Redis.SMembers("newItem").Result()
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Print(services)
	for _, s := range services {
		subscriber := model.Subscriber(s)
		for {
			raw, err := model.Redis.LPop("service-" + s).Result()
			if err != nil {
				break
			}
			var item model.RssItem
			json.Unmarshal([]byte(raw), &item)

			for _, u := range subscriber {
				//bot.PushMessage(u.UserId, linebot.NewTextMessage(item)).Do()
				if _, err := bot.PushMessage(
					u.UserId,
					//linebot.NewFlexMessage(item.Title, template.BubbleContainerRssItem(item)),
					linebot.NewTextMessage(fmt.Sprintf("%s\n%s\n%s\n%s", item.ServiceName, item.Title, item.Published, item.Link)),
				).Do(); err != nil {
					log.Print(err)
				}
			}
		}
		model.Redis.SRem("newItem", s)
	}
}

func cleanRssItems() {
	type Result struct {
		ServiceName string
		Count       int
	}
	var res []Result
	if err := model.DB.Model(model.RssItem{}).Select("service_name, count(*) as count").Group("service_name").Having("count(*) > ?", 100).Scan(&res).Error; err != nil {
		log.Print(err)
	}
	fmt.Println(res)
	for _, r := range res {
		if err := model.DB.Exec("DELETE FROM rss_items WHERE service_name=? AND id < (SELECT id FROM (SELECT id FROM rss_items WHERE service_name=? ORDER BY id DESC LIMIT 1 OFFSET 100) foo)", r.ServiceName, r.ServiceName).Error; err != nil {
			log.Print(err)
		}
	}

}

func ReplyMessage(msg string) error {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(msg),
	).Do(); err != nil {
		return err
	}
	return nil
}

func SourceToID() string {
	switch source.Type {
	case linebot.EventSourceTypeUser:
		return source.UserID
	case linebot.EventSourceTypeGroup:
		return source.GroupID
	case linebot.EventSourceTypeRoom:
		return source.RoomID
	}
	return ""
}

func NotifyWeeklyReport() {
	if _, err := bot.PushMessage(
		MegiId,
		linebot.NewTextMessage("WeeklyReport"),
	).Do(); err != nil {
		log.Print(err)
	}
}
