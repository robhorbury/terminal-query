package utils_test

import (
	"example.com/termquery/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFileQueue(t *testing.T) {
	queue := utils.NewFileQueue()
	assert.Equal(t, queue.Length, 0)
	assert.Nil(t, queue.Head)
	assert.Nil(t, queue.Tail)
}

func TestEnqueue(t *testing.T) {
	queue := utils.NewFileQueue()
	queue.Enqueue("10")
	queue.Enqueue("20")

	assert.Equal(t, queue.Length, 2)
	assert.NotNil(t, queue.Head)
	assert.NotNil(t, queue.Tail)
	assert.NotEqual(t, queue.Head, queue.Tail)
}

func TestDeque(t *testing.T) {
	queue := utils.NewFileQueue()
	queue.Enqueue("10")
	queue.Enqueue("20")

	val, err := queue.Deque()

	assert.Equal(t, queue.Length, 1)
	assert.NotNil(t, queue.Head)
	assert.NotNil(t, queue.Tail)
	assert.Nil(t, err)
	assert.Equal(t, queue.Head, queue.Tail)
	assert.Equal(t, val, "10")
}

func TestRemoveAndDeque(t *testing.T) {
	mockRemoveFunc := func(name string) error { return nil }
	queue := utils.NewFileQueue()
	queue.Enqueue("10")
	queue.Enqueue("20")

	err := queue.RemoveAndDeque("testpath", mockRemoveFunc)

	assert.Equal(t, queue.Length, 1)
	assert.NotNil(t, queue.Head)
	assert.NotNil(t, queue.Tail)
	assert.Nil(t, err)
	assert.Equal(t, queue.Head, queue.Tail)
}
