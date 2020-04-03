package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nattawitc/rich-go/client"
)

const (
	appid = "695217204181860362"
)

func main() {
	defer client.Logout()
	sChan := make(chan state)

	go func() {
		for {
			s, err := getItuneState()
			if err != nil {
				fmt.Println(err)
			}
			select {
			case sChan <- s:
				time.Sleep(time.Second)
			case <-time.After(time.Second):
			}
		}
	}()

	p := &presence{}

	for {
		select {
		case s := <-sChan:
			if !p.oldState.compare(s) {
				p.setRichPresence(s)
			}
		}
	}
}

type presence struct {
	oldState state
	once     sync.Once
	display  bool
}

func (p *presence) setRichPresence(s state) {
	p.once.Do(func() {
		err := client.Login(appid)
		if err != nil {
			panic(err)
		}
	})

	p.oldState = s

	switch s.state {
	case statePlaying:
		if err := setActivity(s.track); err != nil {
			fmt.Println(err)
		}
		p.display = true
	default:
		if p.display {
			p.display = false
			err := client.ClearActivity()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func setActivity(t track) error {
	activity := client.Activity{
		State:      fmt.Sprintf("👤 %v 💿 %v", t.artist, t.album),
		Details:    fmt.Sprintf("🎵 %v", t.name),
		LargeImage: imageName(t.album),
		LargeText:  t.name,
	}

	pp, err := strconv.ParseFloat(t.playerPosition, 64)
	if err != nil {
		return err
	}

	start := time.Now().Add(-time.Duration(pp) * time.Second)

	activity.Timestamps = &client.Timestamps{
		Start: &start,
	}
	fmt.Printf("set activity %+v\n", activity)
	return client.SetActivity(activity)
}

func imageName(album string) string {
	//return album
	album = strings.Replace(album, ":", ";", -1)

	h := sha256.New()
	h.Write([]byte(album + ".jpg"))
	b := h.Sum(nil)
	album = fmt.Sprintf("%x", b[0:10])
	return album
}
