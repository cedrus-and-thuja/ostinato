# Ostinato

A Go library providing functional-style stream processing operations for slices. Ostinato offers a fluent API for chaining operations like map, filter, distinct, reduce, and group-by transformations.

## Installation

```bash
go get github.com/cedrus-and-thuja/ostinato
```

Or add to your `go.mod`:

```
require github.com/cedrus-and-thuja/ostinato v0.1.0
```

## Quick Start

```go
import "github.com/cedrus-and-thuja/ostinato/pkg/ostinato"

// Create a stream from a slice
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
stream := ostinato.Stream(numbers)

// Chain operations
result := stream.
    Filter(func(n int) bool { return n%2 == 0 }).  // Keep even numbers
    Map(func(n int) int { return n * 2 }).         // Double each number
    Distinct(nil)                                   // Remove duplicates

// Convert back to slice (cast to access underlying data)
processed := result.(ostinato.Streamable[int])
```

## Examples

### Basic Transformations

```go
// Map: Transform each element (same type)
numbers := []int{1, 2, 3, 4, 5}
doubled := ostinato.Stream(numbers).
    Map(func(n int) int { return n * 2 })

// MapTo: Transform elements to a different type
numbers := []int{1, 2, 3, 4, 5}
stringified := ostinato.MapTo(ostinato.Stream(numbers), 
    func(n int) string { return fmt.Sprintf("Number: %d", n) })

// Filter: Keep elements matching condition
even := ostinato.Stream(numbers).
    Filter(func(n int) bool { return n%2 == 0 })

// Distinct: Remove duplicates
unique := ostinato.Stream([]int{1, 2, 2, 3, 3, 4}).
    Distinct(nil)
```

### Reduce Operations

```go
// Sum all numbers
numbers := []int{1, 2, 3, 4, 5}
sum := ostinato.Stream(numbers).
    Reduce(func(acc any, n int) any {
        return acc.(int) + n
    }, 0)
// sum = 15

// Concatenate strings
words := []string{"hello", " ", "world"}
sentence := ostinato.Stream(words).
    Reduce(func(acc any, word string) any {
        return acc.(string) + word
    }, "")
// sentence = "hello world"
```

### Group By Operations

```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 45},
    {Name: "Charlie", Age: 35},
    {Name: "David", Age: 55},
}

// Group by age category
grouped := ostinato.Stream(people).
    GroupBy(func(p Person) any {
        if p.Age < 30 {
            return "young"
        } else if p.Age < 50 {
            return "middle"
        }
        return "senior"
    })

// Result: []Grouping[Person, any] with groups by age category
```

### Complex Chains

```go
type Product struct {
    Name  string
    Price float64
    Category string
}

products := []Product{
    {Name: "Laptop", Price: 999.99, Category: "Electronics"},
    {Name: "Book", Price: 15.99, Category: "Books"},
    {Name: "Phone", Price: 699.99, Category: "Electronics"},
    {Name: "Tablet", Price: 399.99, Category: "Electronics"},
}

// Find average price of electronics over $500
expensiveElectronics := ostinato.Stream(products).
    Filter(func(p Product) bool { 
        return p.Category == "Electronics" && p.Price > 500 
    }).
    Reduce(func(acc any, p Product) any {
        data := acc.(map[string]float64)
        data["total"] += p.Price
        data["count"] += 1
        return data
    }, map[string]float64{"total": 0, "count": 0})

avgPrice := expensiveElectronics.(map[string]float64)["total"] / 
           expensiveElectronics.(map[string]float64)["count"]
```

### Type Transformation Example

```go
type Person struct {
    Name string
    Age  int
}

type PersonDTO struct {
    FullName string
    AgeGroup string
}

people := []Person{
    {Name: "Alice Smith", Age: 25},
    {Name: "Bob Johnson", Age: 45},
    {Name: "Charlie Brown", Age: 35},
    {Name: "David Miller", Age: 55},
}

// Transform Person objects to PersonDTO objects
dtos := ostinato.MapTo(ostinato.Stream(people), func(p Person) PersonDTO {
    ageGroup := "adult"
    if p.Age < 30 {
        ageGroup = "young adult"
    } else if p.Age >= 50 {
        ageGroup = "senior"
    }
    
    return PersonDTO{
        FullName: p.Name,
        AgeGroup: ageGroup,
    }
})

// Further process the DTOs
filteredDTOs := dtos.Filter(func(dto PersonDTO) bool {
    return dto.AgeGroup != "adult"
}).ToSlice()
```

## API Reference

### Core Functions

- `Stream[T](slice []T) Streamable[T]` - Create a new stream from a slice
- `Map(fn func(T) T) Streamable[T]` - Transform each element (preserves type)
- `MapTo[T, R](s Streamable[T], fn func(T) R) Streamable[R]` - Transform elements to a different type
- `Filter(fn FilterFunc[T]) Streamable[T]` - Keep elements matching predicate  
- `Distinct(fn IdentityFunc[T, any]) Streamable[T]` - Remove duplicates
- `Reduce(fn func(any, T) any, initial any) any` - Reduce to single value
- `GroupBy(fn func(T) any) any` - Group elements by key function

### Function Types

- `MapFunc[T, R any] func(T) R` - Transformation function
- `FilterFunc[T any] func(T) bool` - Predicate function  
- `IdentityFunc[T any, R any] func(T) R` - Identity/key extraction function
- `ReduceFunc[T, R any] func(R, T) R` - Reduction function

## Testing

Run tests with:

```bash
go test ./pkg/ostinato/...
```

## Requirements

- Go 1.23.5 or later

## License

This project is licensed under the MIT License.
