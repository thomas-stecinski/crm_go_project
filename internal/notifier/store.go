package store

type Notifier interface {
	Send(message string) error
}
