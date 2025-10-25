package service

import (
	"fmt"
	"go-profiler-mem-leak-finder/internal/domain"
	"time"
)

func RunLeakerJob1_SimpleSlice() {
	go func() {
		for {
			data := make([]byte, 1024*1024) // 1MB
			domain.LeakyStore1 = append(domain.LeakyStore1, data)
			time.Sleep(1 * time.Second)
		}
	}()
}

func RunLeakerJob2_ChannelGoroutineBlock() {
	go func() {
		data := make([]byte, 1024*1024) // 1MB
		domain.LeakyChannel2 <- data
	}()
}

func RunLeakerJob3_MutexMap() {
	go func() {
		i := 0
		for {
			data := make([]byte, 1024*1024) // 1MB
			domain.LeakyStore3Mutex.Lock()
			domain.LeakyStore3[i] = data
			domain.LeakyStore3Mutex.Unlock()
			i++
			time.Sleep(1 * time.Second)
		}
	}()
}

func RunSafeJob() {
	go func() {
		fmt.Println("I'm a safe goroutine!")
	}()
}
