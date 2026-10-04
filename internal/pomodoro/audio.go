package pomodoro

import (
	_ "embed"
)

//go:embed data/begin.mp3
var beginMP3 []byte

//go:embed data/break.mp3
var breakMP3 []byte

func playBeginBell() error {
	return playAudio(beginMP3)
}

func playBreakBell() error {
	return playAudio(breakMP3)
}
