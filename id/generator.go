package id

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	timestampBits  = uint(41)
	sequenceBits   = uint64(12)
	workerBits     = uint64(9)
	timestampShift = sequenceBits + workerBits
	sequenceMask   = -1 ^ (-1 << sequenceBits)
	timestampMax   = -1 ^ (-1 << timestampBits)

	high32 = (-1 ^ (-1 << 32)) << 31
)

var (
	gUnixTimestamp int64
)

func init() {
	gUnixTimestamp = getNowTime()
	fmt.Println(gUnixTimestamp)
}

type option func(s *Snowflake1)

func WithWorkerID(workerId int64) option {
	return func(s *Snowflake1) {
		s.workerId = workerId
	}
}

type option2 func(s *Snowflake2)

func WithWorkerID2(workerId int64) option2 {
	return func(s *Snowflake2) {
		s.workerId = workerId
	}
}

type Snowflake1 struct {
	sync.Mutex
	timestamp int64
	sequence  int64
	workerId  int64
}

type Snowflake2 struct {
	sync.Mutex
	workerId  int64
	timestamp int64
}

func NewSnowflake1(op ...option) *Snowflake1 {
	sf := &Snowflake1{}
	for _, o := range op {
		o(sf)
	}
	return sf
}

func (s *Snowflake1) NextID() int64 {
	s.Lock()
	now := getNowTime()
	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = getNowTime()
			}
		}
	} else {
		s.sequence = 0
	}

	t := now - gUnixTimestamp
	if t > timestampMax {
		s.Unlock()
		fmt.Errorf("t must be less than timestamp max")
		return 0
	}
	s.timestamp = now
	r := (t << timestampShift) | (s.sequence << workerBits) | s.workerId
	s.Unlock()
	return r
}

func getNowTime() int64 {
	return time.Now().UnixNano() / 1000000
}

func NewSnowflake2(ops ...option2) *Snowflake2 {
	sf := &Snowflake2{
		workerId:  1,
		timestamp: time.Now().Unix() << 32,
	}
	for _, op := range ops {
		op(sf)
	}

	go func() {
		time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second)))
		tk := time.NewTicker(time.Second)
		now := time.Now().Unix()
		sf.updateTimestamp(now)
		for {
			now := <-tk.C
			sf.updateTimestamp(now.Unix())
		}
	}()

	return sf
}

func (s *Snowflake2) NextID() int64 {
	i := atomic.AddInt64(&s.timestamp, 1)
	return (i & high32) | ((i & sequenceMask) << sequenceBits) | s.workerId
}

func (s *Snowflake2) updateTimestamp(ts int64) {
	s.Lock()
	ts = ts << 32
	if s.timestamp == ts {
		for ts <= s.timestamp {
			ts = time.Now().Unix() << 32
		}
	}
	s.timestamp = ts
	s.Unlock()
}
