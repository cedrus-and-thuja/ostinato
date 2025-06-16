package ostinato

import (
	"fmt"
	"reflect"
	"testing"
)

func cube(v int) int {
	return v * v * v
}

func double(v int) int {
	return v * 2
}

func half(v int) int {
	return v / 2
}

func TestMapSimpleChains(t *testing.T) {
	s := Stream([]int{1, 2, 3})
	s = s.Map(double).Map(half).Map(double)

	expected := []int{2, 4, 6}
	ss := s.(streamable[int])
	if !reflect.DeepEqual(ss.slice, expected) {
		t.Errorf("Expected %v, got %v", expected, ss.slice)
	}
	s = Stream([]int{1, 2, 3})
	s = s.Map(cube)

	expected = []int{1, 8, 27}
	ss = s.(streamable[int])
	if !reflect.DeepEqual(ss.slice, expected) {
		t.Errorf("Expected %v, got %v", expected, ss.slice)
	}
	s = s.Map(half).Map(double)
	expected = []int{0, 8, 26}
	ss = s.(streamable[int])
	if !reflect.DeepEqual(ss.slice, expected) {
		t.Errorf("Expected %v, got %v", expected, ss.slice)
	}
}

func TestMapTo(t *testing.T) {
	s := Stream([]int{1, 2, 3})
	
	// Map integers to strings
	strStream := MapTo(s, func(v int) string {
		return fmt.Sprintf("num-%d", v)
	})
	
	expectedStrs := []string{"num-1", "num-2", "num-3"}
	ss := strStream.(streamable[string])
	if !reflect.DeepEqual(ss.slice, expectedStrs) {
		t.Errorf("Expected %v, got %v", expectedStrs, ss.slice)
	}
	
	// Map integers to a custom struct
	type NumInfo struct {
		Original int
		Squared  int
	}
	
	infoStream := MapTo(s, func(v int) NumInfo {
		return NumInfo{
			Original: v,
			Squared:  v * v,
		}
	})
	
	expectedInfos := []NumInfo{
		{Original: 1, Squared: 1},
		{Original: 2, Squared: 4},
		{Original: 3, Squared: 9},
	}
	
	is := infoStream.(streamable[NumInfo])
	if !reflect.DeepEqual(is.slice, expectedInfos) {
		t.Errorf("Expected %v, got %v", expectedInfos, is.slice)
	}
}

func TestDistinct(t *testing.T) {
	s := Stream([]int{1, 2, 3, 2, 1, 4, 5, 5})
	s = s.Distinct(nil)

	expected := []int{1, 2, 3, 4, 5}
	ss := s.(streamable[int])
	if !reflect.DeepEqual(ss.slice, expected) {
		t.Errorf("Expected %v, got %v", expected, ss.slice)
	}

	s2 := Stream([]string{"a", "b", "a", "c", "b"})
	s2 = s2.Distinct(func(v string) any { return v })

	expectedStr := []string{"a", "b", "c"}
	ss2 := s2.(streamable[string])
	if !reflect.DeepEqual(ss2.slice, expectedStr) {
		t.Errorf("Expected %v, got %v", expectedStr, ss2.slice)
	}
}
func TestDistinctWithIdentity(t *testing.T) {
	s := Stream([]int{1, 2, 3, 2, 1, 4, 5, 5})
	s = s.Distinct(identity)

	expected := []int{1, 2, 3, 4, 5}
	ss := s.(streamable[int])
	if !reflect.DeepEqual(ss.slice, expected) {
		t.Errorf("Expected %v, got %v", expected, ss.slice)
	}
}

func TestDistinctWithIdentityStrings(t *testing.T) {
	s := Stream([]string{"a", "b", "a", "c", "b"})
	s = s.Distinct(identity[string])

	expectedStr := []string{"a", "b", "c"}
	ss := s.(streamable[string])
	if !reflect.DeepEqual(ss.slice, expectedStr) {
		t.Errorf("Expected %v, got %v", expectedStr, ss.slice)
	}
}

func TestDistinctWithNil(t *testing.T) {
	s := Stream([]int{1, 2, 3, 2, 1, 4, 5, 5})
	s = s.Distinct(nil)

	expected := []int{1, 2, 3, 4, 5}
	ss := s.(streamable[int])
	if !reflect.DeepEqual(ss.slice, expected) {
		t.Errorf("Expected %v, got %v", expected, ss.slice)
	}
}

func TestDistinctWithNilStrings(t *testing.T) {
	s := Stream([]string{"a", "b", "a", "c", "b"})
	s = s.Distinct(nil)

	expectedStr := []string{"a", "b", "c"}
	ss := s.(streamable[string])
	if !reflect.DeepEqual(ss.slice, expectedStr) {
		t.Errorf("Expected %v, got %v", expectedStr, ss.slice)
	}
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Equal(p2 Person) bool {
	return p.Name == p2.Name && p.Age == p2.Age
}

func getName(p Person) any {
	return p.Name
}

func TestDistinctWithStruct(t *testing.T) {
	s := Stream([]Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Alice", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Bob", Age: 25},
	})
	s = s.Distinct(getName)

	expected := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	sstream, ok := s.(streamable[Person])
	if !ok {
		t.Fatalf("Expected s to be of type streamable, got %T", s)
	}
	if len(sstream.slice) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, sstream.slice)
	}
	for i, peop := range sstream.slice {
		if !peop.Equal(expected[i]) {
			t.Errorf("Expected %v, got %v", expected[i], peop)
		}
	}
}

func TestFilter(t *testing.T) {
	s := Stream([]int{1, 2, 3, 4, 5, 6})
	s = s.Filter(func(v int) bool {
		return v%2 == 0
	})

	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(s.(streamable[int]).slice, expected) {
		t.Errorf("Expected %v, got %v", expected, s.(streamable[int]).slice)
	}

	s2 := Stream([]string{"apple", "banana", "cherry", "date"})
	s2 = s2.Filter(func(v string) bool {
		return len(v) > 5
	})
	s2Stream, _ := s2.(streamable[string])
	expectedStr := []string{"banana", "cherry"}
	if !reflect.DeepEqual(s2Stream.slice, expectedStr) {
		t.Errorf("Expected %v, got %v", expectedStr, s2Stream.slice)
	}
}

func TestFilterNil(t *testing.T) {
	s := Stream([]int{1, 2, 3, 4, 5, 6})
	s = s.Filter(nil)

	expected := []int{1, 2, 3, 4, 5, 6}
	ss := s.(streamable[int])
	if !reflect.DeepEqual(ss.slice, expected) {
		t.Errorf("Expected %v, got %v", expected, ss.slice)
	}

	s2 := Stream([]string{"apple", "banana", "cherry", "date"})
	s2 = s2.Filter(nil)
	ss2Stream, _ := s2.(streamable[string])
	expectedStr := []string{"apple", "banana", "cherry", "date"}
	if !reflect.DeepEqual(ss2Stream.slice, expectedStr) {
		t.Errorf("Expected %v, got %v", expectedStr, ss2Stream.slice)
	}
}

func TestReduce(t *testing.T) {
	s := Stream([]int{1, 2, 3, 4, 5})
	result := s.Reduce(func(a any, b int) any {
		return a.(int) + b
	}, 0)
	expected := 15
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	s2 := Stream([]string{"a", "b", "c"})
	result2 := s2.Reduce(func(a any, b string) any {
		return a.(string) + b
	}, "")
	expectedStr := "abc"
	if result2 != expectedStr {
		t.Errorf("Expected %v, got %v", expectedStr, result2)
	}
	s = Stream([]int{1, 2, 3, 4, 5})
	result = s.Reduce(func(a any, b int) any {
		return a.(int) * b
	}, 1)
	expectedMul := 120
	if result != expectedMul {
		t.Errorf("Expected %v, got %v", expectedMul, result)
	}
}

func TestReduceWithEmptySlice(t *testing.T) {
	s := Stream([]int{})
	result := s.Reduce(func(a any, b int) any {
		return a.(int) + b
	}, 0)
	expected := 0
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	s2 := Stream([]string{})
	result2 := s2.Reduce(func(a any, b string) any {
		return a.(string) + b
	}, "")
	expectedStr := ""
	if result2 != expectedStr {
		t.Errorf("Expected %v, got %v", expectedStr, result2)
	}
}

func TestReduceOneObjects(t *testing.T) {
	s := Stream([]Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Alice", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Bob", Age: 25},
	})
	s = s.Distinct(getName)
	inital := map[string]int{
		"total": 0,
		"count": 0,
	}
	reduce := func(a any, person Person) any {
		aMap := a.(map[string]int)
		aMap["total"] += person.Age
		aMap["count"]++
		return aMap
	}
	result := s.Reduce(reduce, inital)
	expected := map[string]int{
		"total": 90, // 30 + 25 + 35
		"count": 3,  // Alice, Bob, Charlie
	}
	reduceResult := result.(map[string]int)
	if !reflect.DeepEqual(reduceResult, expected) {
		t.Errorf("Expected %v, got %v", expected, reduceResult)
	}
}
