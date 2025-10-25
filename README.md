# Go Profiler Memory Leak Finder

## 🎯 Project Goal
This project serves as a hands-on demonstration of finding three different types of memory and resource leaks in a Go application that uses goroutines and the Echo web framework. It utilizes the built-in net/http/pprof tools to analyze the application's heap usage and goroutine state to pinpoint the exact source of resource retention.

## 💡 Project Context and Motivation
In modern software development, ensuring efficient memory usage and preventing resource leaks is crucial for building robust applications. Go, with its built-in support for concurrency through goroutines, presents unique challenges in this regard. This project aims to provide developers with practical insights and tools to identify and resolve memory leaks in Go applications, ultimately contributing to better performance and resource management.

Debugging concurrency issues is a cornerstone of professional Go development. This project intentionally simulates three common, yet distinct, patterns of resource leaks to train developers on the proper use of the go tool pprof utility for each scenario:

| Leak Type               | Trigger Endpoint          | Mechanism                                                                                          | Profiling Method                     |
|------------------------|---------------------------|---------------------------------------------------------------------------------------------------|-------------------------------------|
| Leak 1: Simple Slice   | /start-leak/slice         | An endless goroutine continuously appends large byte slices to an unbounded global slice. The global reference prevents the GC from releasing the memory. | Heap Profile (/debug/pprof/heap)    |
| Leak 2: Channel Block   | /start-leak/channel       | A goroutine attempts to send data to an unreceived, buffered channel (size 1). The goroutine blocks indefinitely, resulting in a Goroutine Leak and retention of the data in the channel's buffer/goroutine stack. | Goroutine Profile (/debug/pprof/goroutine) |
| Leak 3: Mutex/Map      | /start-leak/mutex         | An endless goroutine continuously writes large byte slices to an unbounded global map, protected by a sync.Mutex. The map grows forever, leaking memory. | Heap Profile (/debug/pprof/heap)    |

## ✨ Features
The application runs a simple HTTP server (on port 8080) and includes:
- Three Triggerable Leaks: Separate endpoints to start each distinct leak job.
- Safe Goroutine Endpoint (/safe-test): Starts a short-lived goroutine that demonstrates correct resource cleanup.
- Monitoring Endpoint (/leak-test): Reports the current state and size of the active leaks.
- Integrated Profiling: All standard pprof endpoints are registered under /debug/pprof/.

## 🛠️ Prerequisites
- Go (1.18 or higher)
- go tool pprof (comes standard with Go)
- Optional: graphviz (required for generating graphical visualizations using the web command in pprof).

## 🚀 Setup and Run
1. Initialize Go Module (if not already done):
   ```bash
   go mod init go-profiler-mem-leak-finder
   ```
2. Get Dependencies:
   ```bash
   go get github.com/labstack/echo/v4
   go get github.com/labstack/echo-contrib/pprof
   go get -u github.com/swaggo/swag/cmd/swag
   go get -u github.com/swaggo/echo-swagger
    ```
3. Generate Swagger Docs:
   ```bash
   swag init -g cmd/server/main.go
   ```
4. Run the Application:
    ```bash
    go run main.go
    ```
The console will log the server start. No leaks are active yet.

## 🤔 Understanding the Leaks

### Heap Leak

A heap leak occurs when memory is allocated on the heap, but is no longer needed by the application and is not released by the garbage collector. This is typically caused by global variables or other long-lived objects that hold references to objects that are no longer in use.

### Goroutine Leak

A goroutine leak occurs when a goroutine is started but never finishes. This can happen if a goroutine is blocked on a channel that is never written to, or if it is waiting for a condition that never becomes true. Goroutine leaks can lead to increased memory usage and can eventually cause the application to crash.

### Why Each Endpoint is a Leak

*   `/start-leak/slice`: This endpoint starts a goroutine that continuously appends data to a global slice. Because the slice is global, it is never garbage collected, and the memory usage of the application will continue to grow.
*   `/start-leak/channel`: This endpoint starts a goroutine that sends data to a channel. However, there is no other goroutine that is receiving data from the channel. This causes the sending goroutine to block indefinitely, resulting in a goroutine leak.
*   `/start-leak/mutex`: This endpoint starts a goroutine that continuously adds data to a global map. The map is protected by a mutex to prevent race conditions. However, because the map is global, it is never garbage collected, and the memory usage of the application will continue to grow.

## 🚀 How to Use Swagger

Swagger is a tool that allows you to visualize and interact with the API. To use Swagger, navigate to `http://localhost:8080/swagger/index.html` in your browser.

## 🔬 How to Find the Memory Leaks

### Step 1: Trigger the Leaks
Choose the leak(s) you wish to profile. You must hit the endpoint to start the perpetual leak job(s).
```bash
# Example: Start Leak 1 (Slice) AND Leak 3 (Map)
curl http://localhost:8080/start-leak/slice
curl http://localhost:8080/start-leak/mutex
# Optional: Check the current status of the leaks
curl http://localhost:8080/leak-test
```
Wait 30-60 seconds to allow memory to accumulate.
### Step 2: Fetch the Profile
The type of leak determines the profile you should fetch:
| Leak Type(s)          | Profile to Fetch | Command                                               |
|-----------------------|------------------|-------------------------------------------------------|
| Leak 1 or 3 (Memory)  | Heap             | go tool pprof http://localhost:8080/debug/pprof/heap |
| Leak 2 (Goroutine)    | Goroutine        | go tool pprof http://localhost:8080/debug/pprof/goroutine |


### Step 3: Analyze the Profile (Interactive Session)
#### A. Heap Leak Analysis (Leak 1 & 3)
1. Find the Top Allocators: Use the top command in the pprof shell.
   ```bash
   top
   ```
   **Expected Result**: `runLeakerJob1_SimpleSlice` and/or `runLeakerJob3_MutexMap` will show the highest memory usage (inuse_space).
2. Inspect Source Code: Use the list command to see the allocation site.
   ```bash
   list runLeakerJob1_SimpleSlice
   ```
   **Expected Result**: The line with `LeakyStore1 = append(...)` will be highlighted as the primary retention point.
#### B. Goroutine Leak Analysis (Leak 2)
1. View Goroutine Stack Traces: Use the top command on the goroutine profile.
   ```bash
   top
   ```
   **Expected Result**: The top entry will likely be `runLeakerJob2_Channel
GoroutineBlock`, showing a goroutine blocked on a channel operation.
2. Inspect Source Code:
   ```bash
   list runLeakerJob2_ChannelGoroutineBlock
   ```
   **Expected Result**: The line containing `LeakyChannel2 <- data` will be highlighted, showing the goroutine is blocked on a channel send (chan send).
## Available Endpoints for Monitoring
| Endpoint                          | Description                                      |
|-----------------------------------|--------------------------------------------------|
| http://localhost:8080/start-leak/slice   | Starts Leak 1 (Global Slice Retention).          |
| http://localhost:8080/start-leak/channel | Starts Leak 2 (Goroutine Blocking/Channel Leak).   |
| http://localhost:8080/start-leak/mutex   | Starts Leak 3 (Mutex/Map Retention).          |
| http://localhost:8080/leak-test          | Reports status and current size of the active leaks. |
| http://localhost:8080/safe-test          | Triggers a safe, non-leaking goroutine for comparison.      |
| http://localhost:8080/debug/pprof/        | Index page for all profiling endpoints.          |

