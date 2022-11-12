package infra

type ReturningID[T any] struct {
	ID T `ksql:"id"`
}
