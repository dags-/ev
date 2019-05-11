package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/reujab/wallpaper"
	"github.com/getlantern/systray"
	"math/rand"
)

const url = "https://raw.githubusercontent.com/limhenry/earthview/master/earthview.json"

type Wallpaper struct {
	Image string `json:"image"`
}

func main() {
	systray.Run(ready, quit)
}

func ready() {
	var wallpapers []Wallpaper
	fetch(url, &wallpapers)

	interval := time.Minute * 30
	systray.SetTitle("EV")
	next := systray.AddMenuItem("Next", "Skip to the next wallpaper")
	quit := systray.AddMenuItem("Quit", "Stop the app")

	for {
		rand.Seed(time.Now().Unix())
		rand.Shuffle(len(wallpapers), func(i, j int) {
			wallpapers[i], wallpapers[j] = wallpapers[j], wallpapers[i]
		})
		for _, w := range wallpapers {
			wallpaper.SetFromURL(w.Image)
			timer := time.NewTimer(interval)
			select {
			case <- quit.ClickedCh:
				systray.Quit()
				return
			case <-next.ClickedCh:
				continue
			case <-timer.C:
				continue
			}
		}
	}
}

func quit() {

}

func fetch(url string, i interface{}) {
	r, e := http.Get(url)
	if e != nil {
		panic(e)
	}
	defer r.Body.Close()

	e = json.NewDecoder(r.Body).Decode(i)
	if e != nil {
		panic(e)
	}
}
