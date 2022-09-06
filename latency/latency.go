package latency

import (
	"strconv"
	"sync"
	"time"

	cb "github.com/emirpasic/gods/queues/circularbuffer"
	"github.com/metagogs/gogs/gslog"
	"go.uber.org/zap"
)

type msgTime struct {
	agentID int64
	time    int64
}

// LatencyServer 延时管理器
type LatencyServer struct {
	maxQueueCount int //取最近多少次的平均延时
	cache         sync.Map
	latencyQueue  sync.Map
	agentLatency  sync.Map
	// agentCallback map[int64]func(latency int64)
	agentCallback sync.Map
	systemLatency *cb.Queue

	*zap.Logger
}

func NewLatencyServer() *LatencyServer {
	return &LatencyServer{
		maxQueueCount: 10,
		Logger:        gslog.NewLog("latency"),
		systemLatency: cb.New(10),
	}
}

func (l *LatencyServer) OnUpdate(agentID int64, fn func(latency int64)) {
	l.agentCallback.Store(agentID, fn)
}

func (l *LatencyServer) Ping(agentID int64, sendTime int64) {
	key := strconv.FormatInt(agentID, 10) + strconv.FormatInt(sendTime, 10)
	l.cache.Store(key, &msgTime{agentID: agentID, time: sendTime})
	if _, ok := l.latencyQueue.Load(agentID); !ok {
		l.latencyQueue.Store(agentID, cb.New(l.maxQueueCount))
	}
}

func (l *LatencyServer) Pong(agentID int64, revice string) {
	reviceTime := time.Now().UnixMilli()

	key := strconv.FormatInt(agentID, 10) + revice
	if v, ok := l.cache.Load(key); ok { //nolint
		if m, ok := v.(*msgTime); ok {
			l.cache.Delete(key)
			if q, exist := l.latencyQueue.Load(agentID); exist {
				if queue, ok := q.(*cb.Queue); ok {
					queue.Enqueue((reviceTime - m.time) / 2)
					l.getTime(agentID, queue)
				}
			}
		}
	}
}

func (l *LatencyServer) getTime(agentID int64, q *cb.Queue) {
	//算一次延时
	var total int64
	var count int
	for _, t := range q.Values() {
		if v, ok := t.(int64); ok {
			total += v
			count++
		}
	}

	if count > 0 {
		newLatencyTime := total / int64(count)
		l.agentLatency.Store(agentID, newLatencyTime)
		l.Info("get latency",
			zap.Int64("agent_id", agentID),
			zap.Int64("latency", newLatencyTime))

		if fn, ok := l.agentLatency.Load(agentID); ok {
			if fn, ok := fn.(func(latency int64)); ok {
				fn(newLatencyTime)
			}
		}
		l.GetSystemLatency()
	}
}

func (l *LatencyServer) GetSystemLatency() int64 {
	var total int64
	var count int

	l.agentLatency.Range(func(key, value interface{}) bool {
		if v, ok := value.(int64); ok {
			total += v
			count++
		}
		return true
	})

	if count == 0 {
		return 0
	}

	systemTime := total / int64(count)
	l.systemLatency.Enqueue(systemTime)

	return systemTime
}

func (l *LatencyServer) GetSystemLatencyList() []int64 {
	var list []int64
	for _, v := range l.systemLatency.Values() {
		if v, ok := v.(int64); ok {
			list = append(list, v)
		}
	}

	return list
}

func (l *LatencyServer) GetUserLatency() map[int64]int64 {
	var result = make(map[int64]int64)
	l.agentLatency.Range(func(key, value interface{}) bool {
		if v, ok := value.(int64); ok {
			result[key.(int64)] = v
		}
		return true
	})
	return result
}

func (l *LatencyServer) Clear(agentID int64) {
	l.latencyQueue.Delete(agentID)
	l.agentLatency.Delete(agentID)
	l.agentCallback.Delete(agentID)
}
