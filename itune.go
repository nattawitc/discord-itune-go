package main

import (
	"github.com/andybrewer/mack"
)

const (
	statePlaying = "playing"
	statePaused  = "paused"
)

type state struct {
	state string
	track track
}

func (s state) compare(o state) bool {
	return s.state == o.state &&
		s.track.name == o.track.name &&
		s.track.artist == o.track.artist &&
		s.track.album == o.track.album
}

type track struct {
	name           string
	artist         string
	album          string
	playerPosition string // in second
}

func getItuneState() (state, error) {
	s := state{}
	t := &tell{}
	s.state = t.tell("Music", "player state")

	if s.state != statePlaying && s.state != statePaused {
		return s, t.Error()
	}

	s.track.name = t.tell("Music", "name of current track")
	s.track.artist = t.tell("Music", "artist of current track")
	s.track.album = t.tell("Music", "album of current track")
	s.track.playerPosition = t.tell("Music", "player position")

	return s, nil
}

type tell struct {
	err error
}

func (t *tell) tell(application string, commands ...string) string {
	if t.err != nil {
		return ""
	}
	s, err := mack.Tell(application, commands...)
	if err != nil {
		t.err = err
		return ""
	}

	return s
}

func (t *tell) Error() error {
	return t.err
}
