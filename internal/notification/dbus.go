package notification

type Notificator interface {
	Notify(title, message string) error
}

type Pomodoro struct {
	notificator Notificator
}

func New(notificator Notificator) *Pomodoro {
	return &Pomodoro{
		notificator: notificator,
	}
}
