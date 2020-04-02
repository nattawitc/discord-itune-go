package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/nattawitc/rich-go/ipc"
)

const (
	appid = "695217204181860362"
)

func main() {
	defer ipc.CloseSocket()
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
				fmt.Println("state change")
				p.setRichPresence(s)
			}
		}
	}
}

type presence struct {
	oldState state
	once     sync.Once
}

func (p *presence) setRichPresence(s state) {
	fmt.Println("set activity")
	p.once.Do(func() {
		fmt.Println("init")
		err := client.Login(appid)
		if err != nil {
			panic(err)
		}
	})

	p.oldState = s

	switch s.state {
	case statePlaying:
		activity := client.Activity{
			State:      fmt.Sprintf("👤 %v 💿 %v", s.track.artist, s.track.album),
			Details:    fmt.Sprintf("🎵 %v", s.track.name),
			LargeImage: imageName(s.track.album),
			LargeText:  s.track.name,
		}

		pp, err := strconv.ParseFloat(s.track.playerPosition, 64)
		if err != nil {
			fmt.Println(err) // not important error, ignore
		}

		start := time.Now().Add(-time.Duration(pp) * time.Second)

		activity.Timestamps = &client.Timestamps{
			Start: &start,
		}
		fmt.Printf("set activity %+v", activity)
		err = client.SetActivity(activity)
		if err != nil {
			fmt.Println(err)
		}
	}
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
