package cmd

import (
	"rb/internal/notification"
	"rb/internal/pomodoro"

	"github.com/spf13/cobra"
)

var total int

var pomodoroCmd = func() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pomodoro",
		Short: "Starts a pomodoro study session",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			notificator, err := notification.NewDBusNotificator()
			if err != nil {
				return err
			}
			defer notificator.Close()

			p := pomodoro.NewPomodoro(notificator)

			return p.Start(total)
		},
	}

	cmd.Flags().IntVarP(
		&total,
		"total",
		"t",
		75,
		"Total session minutes (study + breaks)",
	)

	return cmd
}()
