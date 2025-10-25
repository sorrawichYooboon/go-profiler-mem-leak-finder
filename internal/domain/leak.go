package domain

import "sync"

var (
	LeakyStore1       [][]byte
	LeakyChannel2     = make(chan []byte, 1)
	LeakyStore3       = make(map[int][]byte)
	LeakyStore3Mutex  = &sync.Mutex{}
)
