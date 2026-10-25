package collection

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

type (
	PriorityQueueSuite struct {
		suite.Suite
		pq Queue[*testPriorityQueueItem]
	}

	testPriorityQueueItem struct {
		value int
	}
)

func testPriorityQueueItemCompareLess(this *testPriorityQueueItem, that *testPriorityQueueItem) bool {
	return this.value < that.value
}

func TestPriorityQueueSuite(t *testing.T) {
	suite.Run(t, new(PriorityQueueSuite))
}

func (s *PriorityQueueSuite) SetupTest() {
	s.pq = NewPriorityQueue(testPriorityQueueItemCompareLess)
}

func (s *PriorityQueueSuite) TestNewPriorityQueueWithItems() {
	items := []*testPriorityQueueItem{
		{value: 10},
		{value: 3},
		{value: 5},
		{value: 4},
		{value: 1},
		{value: 16},
		{value: -10},
	}
	s.pq = NewPriorityQueueWithItems(
		testPriorityQueueItemCompareLess,
		items,
	)

	expected := []int{-10, 1, 3, 4, 5, 10, 16}
	result := []int{}

	for !s.pq.IsEmpty() {
		result = append(result, s.pq.Remove().value)
	}
	s.Equal(expected, result)
}

func (s *PriorityQueueSuite) TestInsertAndPop() {
	s.pq.Add(&testPriorityQueueItem{10})
	s.pq.Add(&testPriorityQueueItem{3})
	s.pq.Add(&testPriorityQueueItem{5})
	s.pq.Add(&testPriorityQueueItem{4})
	s.pq.Add(&testPriorityQueueItem{1})
	s.pq.Add(&testPriorityQueueItem{16})
	s.pq.Add(&testPriorityQueueItem{-10})

	expected := []int{-10, 1, 3, 4, 5, 10, 16}
	result := []int{}

	for !s.pq.IsEmpty() {
		result = append(result, s.pq.Remove().value)
	}
	s.Equal(expected, result)

	s.pq.Add(&testPriorityQueueItem{1000})
	s.pq.Add(&testPriorityQueueItem{1233})
	s.pq.Remove() // remove 1000
	s.pq.Add(&testPriorityQueueItem{4})
	s.pq.Add(&testPriorityQueueItem{18})
	s.pq.Add(&testPriorityQueueItem{192})
	s.pq.Add(&testPriorityQueueItem{255})
	s.pq.Remove() // remove 4
	s.pq.Remove() // remove 18
	s.pq.Add(&testPriorityQueueItem{59})
	s.pq.Add(&testPriorityQueueItem{727})

	expected = []int{59, 192, 255, 727, 1233}
	result = []int{}

	for !s.pq.IsEmpty() {
		result = append(result, s.pq.Remove().value)
	}
	s.Equal(expected, result)
}

func (s *PriorityQueueSuite) TestRemoveClearsPointerFromBackingArrayAndReusesStorage() {
	items := make([]*testPriorityQueueItem, 4, 8)
	for i, value := range []int{1, 2, 3, 4} {
		items[i] = &testPriorityQueueItem{value: value}
	}
	s.pq = NewPriorityQueueWithItems(testPriorityQueueItemCompareLess, items)
	pq := s.pq.(*priorityQueueImpl[*testPriorityQueueItem])

	removedItems := make([]*testPriorityQueueItem, 0, 2)
	for _, expected := range []int{1, 2} {
		removed := s.pq.Remove()
		s.Equal(expected, removed.value)
		removedItems = append(removedItems, removed)
	}

	backingItems := pq.items[:cap(pq.items)]
	for _, removed := range removedItems {
		for _, item := range backingItems {
			s.NotEqual(removed, item, "removed item remains in backing array")
		}
	}

	s.pq.Add(&testPriorityQueueItem{value: 0})
	s.pq.Add(&testPriorityQueueItem{value: 5})

	result := make([]int, 0, s.pq.Len())
	for !s.pq.IsEmpty() {
		result = append(result, s.pq.Remove().value)
	}
	s.Equal([]int{0, 3, 4, 5}, result)
}

func (s *PriorityQueueSuite) TestRandomNumber() {
	for range 1000 {

		expected := []int{}
		result := []int{}
		for range 1000 {
			num := rand.Int()
			s.pq.Add(&testPriorityQueueItem{num})
			expected = append(expected, num)
		}
		sort.Ints(expected)

		for !s.pq.IsEmpty() {
			result = append(result, s.pq.Remove().value)
		}
		s.Equal(expected, result)
	}
}
