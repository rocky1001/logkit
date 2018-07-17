package models

import (
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_KeyValueSlice(t *testing.T) {
	testData := struct {
		origin KeyValueSlice
		expect KeyValueSlice
	}{
		origin: KeyValueSlice{
			{
				"test_start1",
				"",
				"a",
			},
			{
				"test_start2",
				"",
				"a",
			},
			{
				"kafka",
				"",
				"kafa",
			},
			{
				"test_final1",
				"",
				"",
			},
			{
				"test_final2",
				"",
				"",
			},
			{
				"kafkaNew",
				"",
				"kafaNew",
			},
		},
		expect: KeyValueSlice{
			{
				"test_final1",
				"",
				"",
			},
			{
				"test_final2",
				"",
				"",
			},
			{
				"test_start1",
				"",
				"a",
			},
			{
				"test_start2",
				"",
				"a",
			},
			{
				"kafka",
				"",
				"kafa",
			},
			{
				"kafkaNew",
				"",
				"kafaNew",
			},
		},
	}
	sort.Stable(testData.origin)
	assert.Equal(t, testData.expect, testData.origin)
}

func TestErrorQueue(t *testing.T) {
	testErrorQueue := NewErrorQueue(DefaultErrorsListCap)
	assert.Equal(t, DefaultErrorsListCap+1, testErrorQueue.maxSize)
	assert.True(t, testErrorQueue.IsEmpty())
	assert.Equal(t, 0, testErrorQueue.Size())
	assert.Equal(t, 0, testErrorQueue.Front)
	assert.Equal(t, 0, testErrorQueue.Rear)

	testErrorQueue.Put(ErrorInfo{"test error", 123456, 0})
	assert.False(t, testErrorQueue.IsEmpty())
	assert.Equal(t, 1, testErrorQueue.Size())
	errorsList := testErrorQueue.Copy()
	assert.Equal(t, 1, len(errorsList))
	assert.Equal(t, 1, testErrorQueue.ErrorSlice[0].Count)

	for i := 0; i < 80; i++ {
		testErrorQueue.Put(ErrorInfo{"test error", 123456, 0})
	}
	assert.False(t, testErrorQueue.IsEmpty())
	assert.Equal(t, 1, testErrorQueue.Size())
	assert.Equal(t, 0, testErrorQueue.Front)
	assert.Equal(t, 1, testErrorQueue.Rear)
	assert.Equal(t, 81, testErrorQueue.ErrorSlice[0].Count)
	errorsList = testErrorQueue.Copy()
	assert.Equal(t, 1, len(errorsList))

	for i := 0; i < 180; i++ {
		testErrorQueue.Put(ErrorInfo{"test error" + strconv.Itoa(i), 123456, 0})
	}
	assert.False(t, testErrorQueue.IsEmpty())
	assert.Equal(t, DefaultErrorsListCap, testErrorQueue.Size())
	assert.Equal(t, 81, testErrorQueue.Front)
	assert.Equal(t, 80, testErrorQueue.Rear)
	assert.Equal(t, testErrorQueue.Front, testErrorQueue.Rear+1%testErrorQueue.maxSize)
	errorsList = testErrorQueue.Copy()
	assert.Equal(t, DefaultErrorsListCap, len(errorsList))
	assert.Equal(t, 1, testErrorQueue.ErrorSlice[1].Count)

	testErrorQueue.Clear()
	assert.True(t, testErrorQueue.IsEmpty())
	errorsList = testErrorQueue.Copy()
	assert.Equal(t, 0, len(errorsList))
	assert.Equal(t, 0, testErrorQueue.Size())
	assert.Equal(t, 0, testErrorQueue.Front)
	assert.Equal(t, 0, testErrorQueue.Rear)
}
