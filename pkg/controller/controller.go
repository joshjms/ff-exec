package controller

import "sync"

type BoxController struct {
	// BoxManager is a struct that manages the sandbox containers
	Sandboxes []bool
	mu        sync.Mutex
}

var Controller *BoxController

func NewBoxController(n int64) {
	// NewBoxManager creates a new BoxManager with n sandboxes
	Controller = &BoxController{
		Sandboxes: make([]bool, n),
	}
}

func GetEmptyBox() int64 {
	// GetEmptyBox returns the index of an empty sandbox
	Controller.mu.Lock()
	defer Controller.mu.Unlock()

	for i, v := range Controller.Sandboxes {
		if !v {
			Controller.Sandboxes[i] = true
			return int64(i)
		}
	}

	return -1
}

func ReleaseBox(i int64) {
	// ReleaseBox releases the sandbox at index i
	Controller.mu.Lock()
	defer Controller.mu.Unlock()

	Controller.Sandboxes[i] = false
}
