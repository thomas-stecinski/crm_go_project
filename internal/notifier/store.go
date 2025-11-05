// internal/notifier/store.go
package notifier

type Notifier interface {
	Send(message string) error
}
