package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	TIME_SHIFT        = uint(22)
	DATA_CENTER_SHIFT = uint(17)
	WORKER_SHIFT      = uint(12)
	SEQ_SHIFT         = uint(0)

	MAX_DATA_CENTER_NUM = 31
	MAX_WORKER_NUM      = 31
	MAX_SEQ_NUM         = 4095

	TIME_MASK        = int64(4095) << TIME_SHIFT
	DATA_CENTER_MASK = int64(31) << DATA_CENTER_SHIFT
	WORKER_MASK      = int64(31) << WORKER_SHIFT
	SEQ_MASK         = int64(4095) << SEQ_SHIFT

	EPOCH = int64(1483228800000) // 2017-01-01 00:00:00
	// DEFAULT_DATA_CENTER = 0
	// DEFAULT_WORKER      = 0
)

type Leaf struct {
	lastTimestamp int64
	dataCenterId  int64
	workerId      int64
	seq           int64
	lock          sync.Mutex
}

func NewLeaf(dataCenterId, workerId int64) (*Leaf, error) {
	if dataCenterId < 0 || dataCenterId > MAX_DATA_CENTER_NUM {
		return nil, errors.New("invalid data center id")
	}
	if workerId < 0 || workerId > MAX_WORKER_NUM {
		return nil, errors.New("invalid worker id")
	}
	return &Leaf{
		dataCenterId: dataCenterId,
		workerId:     workerId,
	}, nil
}

func (l *Leaf) NextId() (int64, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / 1000000

	if now < l.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	if now == l.lastTimestamp {
		l.seq = (l.seq + 1) & MAX_SEQ_NUM
		if l.seq == 0 {
			now = l.waitNextMillis(now)
		}
	} else {
		l.seq = 0
	}

	l.lastTimestamp = now

	return ((now - EPOCH) & TIME_MASK) |
		(l.dataCenterId << DATA_CENTER_SHIFT) |
		(l.workerId << WORKER_SHIFT) |
		(l.seq << SEQ_SHIFT), nil
}

func (l *Leaf) waitNextMillis(now int64) int64 {
	for now <= l.lastTimestamp {
		now = time.Now().UnixNano() / 1000000
	}
	return now
}

// func main() {
// 	leaf, err := NewLeaf(DEFAULT_DATA_CENTER, DEFAULT_WORKER)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i := 0; i < 10; i++ {
// 		id, err := leaf.NextId()
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(id)
// 	}
// }
