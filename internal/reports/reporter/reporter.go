package reporter

type Publisher[T any] interface {
	Register(chan T)
	Unregister(chan T)
}
