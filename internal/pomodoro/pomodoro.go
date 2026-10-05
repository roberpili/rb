package pomodoro

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"rb/internal/notification"
	"syscall"
	"time"
)

type Pomodoro struct {
	notificator notification.Notificator
}

func NewPomodoro(notificator notification.Notificator) *Pomodoro {
	return &Pomodoro{
		notificator: notificator,
	}
}

func (p *Pomodoro) Start(totalMinutes int) error {
	if totalMinutes <= 0 {
		return fmt.Errorf("total must be > 0")
	}

	// allow ctrl+c
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// build the session
	plan := buildPlan(totalMinutes)

	rounds := 0
	for _, s := range plan {
		if s.label == "STUDY" {
			rounds++
		}
	}
	fmt.Printf("Starting %d minutes pomodoro session (%d blocks)\n\n", totalMinutes, rounds)

	for _, seg := range plan {
		if seg.label == "STUDY" {
			if err := playBeginBell(); err != nil {
				return fmt.Errorf("failed to play study bell: %w", err)
			}
		} else {
			if err := playBreakBell(); err != nil {
				return fmt.Errorf("failed to play break bell: %w", err)
			}
		}

		if err := runTimer(ctx, seg.minutes, seg.label); err != nil {
			return err
		}

		fmt.Printf("\n%s COMPLETED\n\n", seg.label)

		if err := p.notificator.Notify(
			seg.label+" finished",
			"Time for the next block",
		); err != nil {
			fmt.Fprintf(os.Stderr, "notification failed: %v\n", err)
		}
	}

	// finished
	if err := playBeginBell(); err != nil {
		return fmt.Errorf("failed to play finish bell: %w", err)
	}

	if err := p.notificator.Notify("Pomodoro finished", "Study session is completed"); err != nil {
		fmt.Fprintf(os.Stderr, "notification failed: %v\n", err)
	}

	return nil
}

// replace time.Sleep allowing context
func runTimer(ctx context.Context, minutes int, label string) error {
	fmt.Printf("[%s] Started: %d minutes\n", label, minutes)

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for consumed := 1; consumed <= minutes; consumed++ {
		select {
		case <-ctx.Done():
			fmt.Println("\n\ncancelled by ctrl+c.")
			return ctx.Err()

		case <-ticker.C:
			remaining := minutes - consumed
			fmt.Printf("[%s] Consumed: %d minutes | Remaining: %d minutes\n", label, consumed, remaining)
		}
	}

	return nil
}
