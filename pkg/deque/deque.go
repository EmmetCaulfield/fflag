package deque

type Deque[T comparable] []T

type Empty struct{}
func (e *Empty) Error() string {
	return "empty deque"
}

func (d *Deque[T]) Push(item T) {
	*d = append(*d, item)
}

func (d *Deque[T]) Append(items ...T) {
	*d = append(*d, items...)
}

func (d *Deque[T]) Pop() (T, error) {
	if len(*d) == 0 {
		return *new(T), &Empty{}
	}
	top := (*d)[len(*d)-1]
	*d = (*d)[0:len(*d)-1]
	return top, nil
}

func (d *Deque[T]) Peek() (T, error) {
	if len(*d) == 0 {
		return *new(T), &Empty{}
	}
	return (*d)[len(*d)-1], nil
}

func (d *Deque[T]) Unshift(item T) {
	*d = append([]T{item}, *d...)
}

func (d *Deque[T]) Prepend(items ...T) {
	*d = append(items, *d...)
}

func (d *Deque[T]) Shift() (T, error) {
	if len(*d) == 0 {
		return *new(T), &Empty{}
	}
	front := (*d)[0]
	*d = (*d)[1:len(*d)]
	return front, nil
}

func (d *Deque[T]) Front() (T, error) {
	if len(*d) == 0 {
		return *new(T), &Empty{}
	}
	return (*d)[0], nil
}

func (d *Deque[T]) Clear() {
	*d = (*d)[:0]
}

func (d *Deque[T]) Init(items ...T) {
	d.Clear()
	d.Append(items...)
}

func EqualV[T comparable](a, b Deque[T]) bool {
	if len(a) != len(b) {
		return false
	}
	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func EqualP[T comparable](a, b *Deque[T]) bool {
	if len(*a) != len(*b) {
		return false
	}
	for i, v := range *a {
		if v != (*b)[i] {
			return false
		}
	}
	return true
}

func (d *Deque[T]) Equal(a *Deque[T]) bool {
	return EqualP(d, a)
}
