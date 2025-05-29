package utils

import (
	"fmt"
	"path/filepath"
)

type qnode struct {
	path string
	next *qnode
}

type FileQueue struct {
	Length int
	Head   *qnode
	Tail   *qnode
}

func NewFileQueue() *FileQueue {
	return &FileQueue{Length: 0, Head: nil, Tail: nil}
}

func (queue *FileQueue) Enqueue(path string) {
	newNode := qnode{path: path, next: nil}
	if queue.Head == nil && queue.Tail == nil {
		queue.Head = &newNode
		queue.Tail = &newNode
	} else {
		queue.Tail.next = &newNode
		queue.Tail = &newNode
	}
	queue.Length = queue.Length + 1
}

func (queue *FileQueue) Deque() (string, error) {
	if queue.Length == 0 {
		return "", fmt.Errorf("queue is empty")
	} else if queue.Head == nil && queue.Tail == nil {
		return "", fmt.Errorf("queue is empty")
	} else if queue.Length == 1 {
		current := queue.Head
		queue.Head = nil
		queue.Tail = nil
		queue.Length = queue.Length - 1
		current.next = nil
		return current.path, nil
	} else {
		current := queue.Head
		queue.Head = queue.Head.next
		queue.Length = queue.Length - 1
		current.next = nil
		return current.path, nil
	}
}

func (queue *FileQueue) Peak() (string, error) {
	if queue.Length == 0 {
		return "", fmt.Errorf("queue is empty")
	} else if queue.Head == nil && queue.Tail == nil {
		return "", fmt.Errorf("queue is empty")
	} else if queue.Length == 1 {
		return queue.Head.path, nil
	} else {
		return queue.Head.path, nil
	}
}

func (queue *FileQueue) RemoveAndDeque(cachePath string, removeFunc RemoveFunc) error {
	fileName, err := queue.Deque()
	if err != nil {
		return err
	}
	err = removeFunc(filepath.Join(cachePath, fileName))
	if err != nil {
		return err
	}
	return nil
}
