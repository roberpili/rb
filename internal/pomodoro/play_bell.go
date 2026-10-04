package pomodoro

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

var (
	otoCtx      *oto.Context
	initOtoOnce sync.Once
	initOtoErr  error
)

// oto needs to init once. this will use your def audio device
// rate 44100 or 48000. other values might cause distortions
// channel count 1 is mono sound, and 2 is stereo
func getOtoContext() (*oto.Context, error) {
	initOtoOnce.Do(func() {
		op := &oto.NewContextOptions{
			SampleRate:   44100,
			ChannelCount: 2,
			Format:       oto.FormatSignedInt16LE,
		}

		ctx, readyChan, err := oto.NewContext(op)
		if err != nil {
			initOtoErr = fmt.Errorf("oto.NewContext failed: %w", err)
			return
		}

		// it might take a bit for the hardware audio devices to be ready, so we wait on the channel
		<-readyChan
		if err := ctx.Err(); err != nil {
			initOtoErr = fmt.Errorf("oto initialization failed: %w", err)
			return
		}

		otoCtx = ctx
	})

	return otoCtx, initOtoErr
}

func playAudio(data []byte) error {
	ctx, err := getOtoContext()
	if err != nil {
		return err
	}

	decodedMp3, err := mp3.NewDecoder(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("mp3.NewDecoder failed: %w", err)
	}

	player := ctx.NewPlayer(decodedMp3)
	player.Play()

	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
