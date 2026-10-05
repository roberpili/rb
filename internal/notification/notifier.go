package notification

import "github.com/godbus/dbus/v5"

type DBusNotificator struct {
	conn *dbus.Conn
}

func NewDBusNotificator() (*DBusNotificator, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	return &DBusNotificator{conn: conn}, nil
}

func (n *DBusNotificator) Notify(title, message string) error {

	obj := n.conn.Object(
		"org.freedesktop.Notifications",
		"/org/freedesktop/Notifications",
	)

	return obj.Call(
		"org.freedesktop.Notifications.Notify",
		0,
		"", uint32(0),
		"",
		title,
		message,
		[]string{},
		map[string]dbus.Variant{},
		int32(5000),
	).Err
}

func (n *DBusNotificator) Close() {
	n.conn.Close()
}
