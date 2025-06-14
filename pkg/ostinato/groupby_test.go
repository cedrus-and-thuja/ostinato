package ostinato

import (
	"reflect"
	"testing"
)

func TestGroupByNumerics(t *testing.T) {
	s := Stream([]int{1, 2, 3, 4, 5, 6})
	grouped := s.GroupBy(func(v int) any { return v % 2 })

	expected := []Grouping[int, any]{
		{
			Key:    1,
			Values: []int{1, 3, 5},
		},
		{
			Key:    0,
			Values: []int{2, 4, 6},
		},
	}
	reGrouped := grouped.([]Grouping[int, any])
	// Since map iteration order is not guaranteed, we need to check the contents
	// without relying on the order
	if len(reGrouped) != len(expected) {
		t.Errorf("Expected %d groups, got %d", len(expected), len(reGrouped))
		return
	}
	
	// Create maps for easy lookup
	actualByKey := make(map[int][]int)
	for _, g := range reGrouped {
		key := g.Key.(int)
		actualByKey[key] = g.Values
	}
	
	expectedByKey := make(map[int][]int)
	for _, g := range expected {
		key := g.Key.(int)
		expectedByKey[key] = g.Values
	}
	
	// Compare the groups
	for key, expectedValues := range expectedByKey {
		actualValues, ok := actualByKey[key]
		if !ok {
			t.Errorf("Missing group for key %v", key)
			continue
		}
		
		if !reflect.DeepEqual(actualValues, expectedValues) {
			t.Errorf("For key %v, expected %v, got %v", key, expectedValues, actualValues)
		}
	}
}

func TestGroupByPeople(t *testing.T) {
	s := Stream([]Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 45},
		{Name: "Charlie", Age: 53},
		{Name: "David", Age: 60},
		{Name: "Eve", Age: 25},
	})
	grouped := s.GroupBy(func(v Person) any {
		if v.Age < 40 {
			return "young"
		} else if v.Age < 60 {
			return "middle-aged"
		}
		return "senior"
	})

	expected := []Grouping[Person, any]{
		{
			Key: "young",
			Values: []Person{
				{Name: "Alice", Age: 30},
				{Name: "Eve", Age: 25},
			},
		},
		{
			Key: "middle-aged",
			Values: []Person{
				{Name: "Bob", Age: 45},
				{Name: "Charlie", Age: 53},
			},
		},
		{
			Key: "senior",
			Values: []Person{
				{Name: "David", Age: 60},
			},
		},
	}
	reGrouped := grouped.([]Grouping[Person, any])
	// Since map iteration order is not guaranteed, we need to check the contents
	// without relying on the order
	if len(reGrouped) != len(expected) {
		t.Errorf("Expected %d groups, got %d", len(expected), len(reGrouped))
		return
	}
	
	// Create maps for easy lookup
	actualByKey := make(map[string][]Person)
	for _, g := range reGrouped {
		key := g.Key.(string)
		actualByKey[key] = g.Values
	}
	
	expectedByKey := make(map[string][]Person)
	for _, g := range expected {
		key := g.Key.(string)
		expectedByKey[key] = g.Values
	}
	
	// Compare the groups
	for key, expectedValues := range expectedByKey {
		actualValues, ok := actualByKey[key]
		if !ok {
			t.Errorf("Missing group for key %v", key)
			continue
		}
		
		if !reflect.DeepEqual(actualValues, expectedValues) {
			t.Errorf("For key %v, expected %v, got %v", key, expectedValues, actualValues)
		}
	}
}
