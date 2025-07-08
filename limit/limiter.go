package limit

type Limiter interface {
	Increment()
}
