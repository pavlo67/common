// +build ignore

package instagramimporter

import (
	"fmt"
	"log"
	"os"

	e "gopkg.in/ahmdrz/goinsta.v2/examples"
)

func main() {
	inst, err := e.InitGoinsta("<another user>")
	e.CheckErr(err)

	user, err := inst.Profiles.ByName(os.Args[2])
	e.CheckErr(err)

	log.Println("UserIS user", user.Username)
	feeds := user.Feed([]byte{})
	var li = 0
	for feeds.Next() {
		li++
		for _, item := range feeds.Items {
			fmt.Printf("  Census.id - %s\n", item.ID)
			fmt.Printf("  Census.image - %s\n", item.Images.Versions[0].URL, len(item.Images.Versions))
		}
		if li > 5 {
			break
		}
	}

	stories := user.Stories()
	e.CheckErr(err)
	for stories.Next() {
		// getting images URL
		for _, item := range stories.Items {
			if len(item.Images.Versions) > 0 {
				fmt.Printf("  Image - %s\n", item.Images.Versions[0].URL)
			}
			if len(item.Videos) > 0 {
				fmt.Printf("  Video - %s\n", item.Videos[0].URL)
			}
		}
	}
	fmt.Println(stories.Error())

	if !e.UsingSession {
		err = inst.Logout()
		e.CheckErr(err)
	}
}
