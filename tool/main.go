package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"os/exec"
)

func main() {
	/*
		go cmd.Manager.Work()
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			done := cmd.Manager.SendCmd(scanner.Text())
			fmt.Println(<-done)
		}
	*/

	flagServiceUrl := "http://yepan.net/bbs/rss.php?bo_table=comm_info"
	out, err := exec.Command("curl", "-c", "cookie", "-XGET", "-L", flagServiceUrl).Output()
	if err != nil {
		log.Print("curl:" + err.Error())
		return
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(out))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed.Title, feed.Description)
}
