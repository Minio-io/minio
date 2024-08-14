package once

// Singleton contains a pointer to T that must be set once.
// Until the value is set all Get() calls will block.
type Singleton[T any] struct {
	v   *T
	set chan struct{}
}

// NewSingleton creates a new unset singleton.
func NewSingleton[T any]() *Singleton[T] {
	return &Singleton[T]{set: make(chan struct{}), v: nil}
}

// Get will return the singleton value.
func (s *Singleton[T]) Get() *T {
	<-s.set
	return s.v
}

// Set the value and unblock all Get requests.
// This may only be called once, a second call will panic.
func (s *Singleton[T]) Set(v *T) {
	s.v = v
	close(s.set)
}
