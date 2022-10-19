package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/e2e"
	"github.com/metagogs/gogs/e2e/testdata"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/logic/baseworld"
	"github.com/metagogs/gogs/e2e/testdata/game"
	"github.com/metagogs/gogs/global"
)

// go tool pprof -http=:8081  mem.pprof
// go tool pprof -http=:8081  cpu.pprof
// go tool trace trace.out
//
// GODEBUG=gctrace=1  go run main.go
//
// gc 5 @16.016s 0%: 0.069+0.62+0.004 ms clock, 0.83+0.34/1.5/0.72+0.053 ms cpu, 6->7->2 MB, 6 MB goal, 0 MB stacks, 0 MB globals, 12 P
// gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
// where the fields are as follows:
//
//	gc #         the GC number, incremented at each GC
//	@#s          time in seconds since program start
//	#%           percentage of time spent in GC since program start
//	#+...+#      wall-clock/CPU times for the phases of the GC,
//
// stop-the-world (STW) sweep termination, concurrent mark and scan, and STW mark termination
//
//	#->#-># MB   heap size at GC start, at GC end, and live heap
//	# MB goal    goal heap size
//	# MB stacks  estimated scannable stack size
//	# MB globals scannable global size
//	# P          number of processors used
func main() {
	startTime := time.Now()
	// debug.SetMemoryLimit(1 * 1024 * 1024 * 1024) // 1GB
	// debug.SetGCPercent(-1)
	global.GOGS_DISABLE_LOG = true
	timeEnd := 30
	cliens := 50

	f, _ := os.Create("cpu.pprof")
	defer f.Close()

	fr, _ := os.Create("trace.out")
	defer fr.Close()
	_ = pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	_ = trace.Start(fr)
	defer trace.Stop()

	serverConfig := &config.Config{}
	serverConfig.Debug = false
	serverConfig.AgentHeartBeatTimeout = 300

	cancel, started := startServer(serverConfig)
	// use default encode proto
	<-started
	testdata.TestApp.UseDefaultEncodeProto()
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	initGC := m.NumGC

	clients := []*e2e.TestClient{}
	for i := 0; i < cliens; i++ {
		client, _ := e2e.NewWSClinet("ws://127.0.0.1:8888/base")
		go client.Start2()

		_ = client.WriteMessage(websocket.BinaryMessage, encodeMessage(&game.BindUser{
			Uid: fmt.Sprintf("test_%d", i),
		}))
		<-baseworld.BindUserHandler
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("user join timeout")
		case <-client.Datas:
			_ = client.WriteMessage(websocket.BinaryMessage, encodeMessage(&game.JoinWorld{
				Uid: strconv.Itoa(i),
			}))
			<-baseworld.JoinWorldHandler
			clients = append(clients, client)
		}

	}

	for _, client := range clients {
		go func(c *e2e.TestClient) {
			run(c)
		}(client)
	}

	<-time.After(time.Duration(timeEnd) * time.Second)
	fmt.Printf("test end, cost time %d \n", time.Since(startTime).Milliseconds())

	fm, _ := os.Create("mem.pprof")
	defer fm.Close()
	_ = pprof.WriteHeapProfile(fm)

	cancel()
	//print system info
	m = runtime.MemStats{}
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc %d\n", m.Alloc)
	fmt.Printf("TotalAlloc %d\n", m.TotalAlloc)
	fmt.Printf("Mallocs %d\n", m.Mallocs)
	fmt.Printf("Frees %d\n", m.Frees)
	fmt.Printf("HeapAlloc %d\n", m.HeapAlloc)
	fmt.Printf("NumGC %d\n", m.NumGC)
	fmt.Printf("InitNumGC %d\n", initGC)
	fmt.Printf("NumForcedGC %d\n", m.NumForcedGC)
	fmt.Printf("PauseTotalNs %dms\n", m.PauseTotalNs/1000/1000)
}

func run(client *e2e.TestClient) {
	ticker := time.NewTicker(16 * time.Millisecond)
	messageData := encodeMessage(&game.UpdateUserInWorld{
		Uid: "1",
		Position: &game.Vecotr3{
			X: 11,
			Y: 22,
			Z: 33,
		},
	})
	for range ticker.C {
		_ = client.WriteMessage(websocket.BinaryMessage, messageData)
	}
}

func startServer(config *config.Config) (context.CancelFunc, chan struct{}) { //nolint
	startTest := make(chan struct{})
	config.StaredCallback = func() {
		<-time.After(1 * time.Second)
		startTest <- struct{}{}
	}
	ctx, cancel := context.WithCancel(context.Background())
	go testdata.StartServer(ctx, config)
	return cancel, startTest
}

func encodeMessage(in interface{}) []byte { //nolint
	p, _ := testdata.TestApp.MessageServer.EncodeMessage(in)
	return p.ToData().B
}
