package lib

type Queue[T any] struct {
	data []T
}

func NewQueue[T any]() Queue[T] {
	return Queue[T]{make([]T, 0)}
}

func (q *Queue[T]) Enqueue(p T) {
	q.data = append(q.data, p)
}

func (q *Queue[T]) Dequeue() T {
	e := q.data[0]
	q.data = q.data[1:]
	return e
}

func (q Queue[T]) IsEmpty() bool {
	return len(q.data) == 0
}
