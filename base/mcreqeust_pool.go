package base

import (
	"github.com/couchbase/gomemcached"
	"github.com/couchbase/goxdcr/log"
	"sync"
)

type MCRequestPool struct {
	name     string
	obj_pool *sync.Pool
	lock     *sync.RWMutex
	logger *log.CommonLogger
}

func NewMCRequestPool(name string, logger *log.CommonLogger) *MCRequestPool {
	return &MCRequestPool{name: name,
		obj_pool: &sync.Pool{},
		lock:     &sync.RWMutex{},
		logger: logger,
	}
}

func (pool *MCRequestPool) Get() *WrappedMCRequest {
	var obj_ret *WrappedMCRequest = nil
	obj := pool.obj_pool.Get()
	if obj == nil {
		obj = pool.addOne()
	}
	obj_ret, ok := obj.(*WrappedMCRequest)
	if !ok {
		panic("object in MCRequestPool should be of type *WrappedMCRequest")
	}

	return obj_ret
}

func (pool *MCRequestPool) addOne() *WrappedMCRequest {
	obj := &WrappedMCRequest{Seqno: 0,
		Req: &gomemcached.MCRequest{Extras: make([]byte, 24)},
	}

	return obj
}

func (pool *MCRequestPool) Put(req *WrappedMCRequest) {
	//make the request vanilar
	req_clean := pool.cleanReq(req)
	pool.obj_pool.Put(req_clean)
}

func (pool *MCRequestPool) cleanReq(req *WrappedMCRequest) *WrappedMCRequest {
	req.Req = pool.cleanMCReq(req.Req)
	req.Seqno = 0
	return req
}

func (pool *MCRequestPool) cleanMCReq(req *gomemcached.MCRequest) *gomemcached.MCRequest {
	req.Cas = 0
	req.Opaque = 0
	req.VBucket = 0
	req.Key = nil
	req.Body = nil
	pool.cleanExtras(req)
	//opCode
	req.Opcode = 0

	return req
}

func (pool *MCRequestPool) cleanExtras(req *gomemcached.MCRequest) {
	for i := 0; i < 24; i++ {
		req.Extras[0] = 0
	}
}
