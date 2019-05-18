package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/getlantern/systray"
	"github.com/reujab/wallpaper"
)

const url = "https://raw.githubusercontent.com/limhenry/earthview/master/earthview.json"

type Wallpaper struct {
	Image string `json:"image"`
}

type Config struct {
	Current string `json:"current"`
	Pause   bool   `json:"pause"`
}

var (
	box = rice.MustFindBox("assets")
)

func main() {
	systray.Run(ready, func() {})
}

func ready() {
	var wallpapers []Wallpaper
	fetch(url, &wallpapers)

	interval := time.Minute * 60
	systray.SetIcon(icon())
	pause := systray.AddMenuItem("Pause", "Keep the current wallpaper")
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
			for first, paused := true, false; first || paused; first = false {
				select {
				case <-quit.ClickedCh:
					systray.Quit()
					return
				case <-next.ClickedCh:
					paused = false
					togglePause(pause, paused)
					continue
				case <-pause.ClickedCh:
					paused = !paused
					togglePause(pause, paused)
					continue
				case <-timer.C:
					continue
				}
			}
		}
	}
}

func togglePause(i *systray.MenuItem, pause bool) {
	if pause {
		i.SetTitle("Unpause")
	} else {
		i.SetTitle("Pause")
	}
}

func icon() []byte {
	if runtime.GOOS == "windows" {
		return box.MustBytes("icon.ico")
	}
	return box.MustBytes("icon.png")
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
