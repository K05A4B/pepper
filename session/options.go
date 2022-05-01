package session

type Options struct {
	LifeCycle     int64
	CleanInterval int64
	HandlerGCError func(error)
}