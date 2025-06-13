package ostinato

// takes an object and returns a transformation of it
type MapFunc[T, R any] func(T) R

// filters the stream based on a condition
// returns true if the object should be included in the stream
type FilterFunc[T any] func(T) bool

// function that takes an object and returns a unique identifier for it
// used for distinct operations
// if nil is passed, the identity function is used
// identity function returns the object itself
type IdentityFunc[T any, R any] func(T) R

// retuces objects in the stream to an accumulator
// takes the object and the accumulator and returns the accumulator (or a new accumulator)
type ReduceFunc[T, R any] func(R, T) R

type GroupByFunc[T any, K any] func(T) K
type Grouping[T any, K any] struct {
	Key    K
	Values []T
}

type Streamable[T any] interface {
	Map(MapFunc[T, T]) Streamable[T]
	Distinct(IdentityFunc[T, any]) Streamable[T]
	Filter(FilterFunc[T]) Streamable[T]

	// Reduce applies a reduction function to the stream, returning a single value.
	// The initial value can be provided, or if nil is passed, the first element of the stream is used as the initial value.
	// If the stream is empty and no initial value is provided, it returns nil.
	// ReduceFunc takes two parameters: the current accumulated value and the next value from the stream.
	// It returns the accumulated value after processing all elements.
	Reduce(func(any, T) any, any) any

	GroupBy(func(T) any) any
}

func identity[T any](v T) any {
	return v
}

func Stream[T any](slice []T) Streamable[T] {
	return streamable[T]{slice: slice}
}

type streamable[T any] struct {
	slice []T
}

func (s streamable[T]) Map(fn MapFunc[T, T]) Streamable[T] {
	newSlice := make([]T, len(s.slice))
	for i, v := range s.slice {
		newSlice[i] = fn(v)
	}
	s.slice = newSlice
	return s
}

func (s streamable[T]) Distinct(fn IdentityFunc[T, any]) Streamable[T] {
	if fn == nil {
		return s.Distinct(identity[T])
	}
	seen := make(map[any]bool)
	newSlice := []T{}
	for _, v := range s.slice {
		id := fn(v)
		if !seen[id] {
			seen[id] = true
			newSlice = append(newSlice, v)
		}
	}
	s.slice = newSlice
	return s
}

func (s streamable[T]) Filter(fn FilterFunc[T]) Streamable[T] {
	if fn == nil {
		return s
	}
	newSlice := []T{}
	for _, v := range s.slice {
		if fn(v) {
			newSlice = append(newSlice, v)
		}
	}
	s.slice = newSlice
	return s
}

func (s streamable[T]) Reduce(fn func(any, T) any, initial any) any {
	result := initial
	if result == nil && len(s.slice) > 0 {
		result = s.slice[0]
		s.slice = s.slice[1:] // Remove the first element since it's used as the initial value
	}
	if fn == nil {
		fn = func(acc any, v T) any {
			return v
		}
	}
	if len(s.slice) == 0 && result == nil {
		return nil // If the stream is empty and no initial value is provided, return nil
	}
	for _, v := range s.slice {
		result = fn(result, v)
	}
	return result
}

func (s streamable[T]) GroupBy(fn func(T) any) any {
	if fn == nil {
		fn = identity[T]
	}
	grouped := make(map[any][]T)
	for _, v := range s.slice {
		key := fn(v)
		grouped[key] = append(grouped[key], v)
	}
	result := make([]Grouping[T, any], 0, len(grouped))
	for key, values := range grouped {
		nvals := values
		result = append(result, Grouping[T, any]{
			Key:    key,
			Values: nvals,
		})
	}
	return result
}
