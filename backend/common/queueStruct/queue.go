package queueStruct

import (
	"reflect"
	"slices"
)

type Queue struct {
	Data []any
}

// Enqueue method
func (q *Queue) Enqueue(element any) {
	q.Data = append(q.Data, element)
}

// Dequeue method
func (q *Queue) Dequeue() (any, bool) {
	if len(q.Data) == 0 {
		return 0, false
	}
	element := q.Data[0]
	q.Data = q.Data[1:]
	return element, true
}

func (q *Queue) Remove(element any) bool {
	for i, ele := range q.Data {
		if reflect.DeepEqual(ele, element) {
			q.Data = slices.Delete(q.Data, i, i+1)
			return true
		}
	}
	return false
}
