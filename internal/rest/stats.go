package rest

import "sync/atomic"

type Stats struct {
	restRequests      uint64
	restReads         uint64
	restIngest        uint64
	restRejected      uint64
	restUnauthorized  uint64
	ingestBatches     uint64
	ingestWrittenRegs uint64
	ingestRejected    uint64
}

func (s *Stats) IncRequests()     { atomic.AddUint64(&s.restRequests, 1) }
func (s *Stats) IncReads()        { atomic.AddUint64(&s.restReads, 1) }
func (s *Stats) IncIngest()       { atomic.AddUint64(&s.restIngest, 1) }
func (s *Stats) IncRejected()     { atomic.AddUint64(&s.restRejected, 1) }
func (s *Stats) IncUnauthorized() { atomic.AddUint64(&s.restUnauthorized, 1) }

func (s *Stats) IncIngestBatch() { atomic.AddUint64(&s.ingestBatches, 1) }
func (s *Stats) AddWrittenRegs(n uint32) {
	atomic.AddUint64(&s.ingestWrittenRegs, uint64(n))
}
func (s *Stats) IncIngestRejected() { atomic.AddUint64(&s.ingestRejected, 1) }

type StatsSnapshot struct {
	REST struct {
		Requests     uint64 `json:"requests"`
		Reads        uint64 `json:"reads"`
		Ingest       uint64 `json:"ingest"`
		Rejected     uint64 `json:"rejected"`
		Unauthorized uint64 `json:"unauthorized"`
	} `json:"rest"`
	Ingest struct {
		Batches  uint64 `json:"batches"`
		Written  uint64 `json:"written"`
		Rejected uint64 `json:"rejected"`
	} `json:"ingest"`
}

func (s *Stats) Snapshot() StatsSnapshot {
	var out StatsSnapshot
	out.REST.Requests = atomic.LoadUint64(&s.restRequests)
	out.REST.Reads = atomic.LoadUint64(&s.restReads)
	out.REST.Ingest = atomic.LoadUint64(&s.restIngest)
	out.REST.Rejected = atomic.LoadUint64(&s.restRejected)
	out.REST.Unauthorized = atomic.LoadUint64(&s.restUnauthorized)

	out.Ingest.Batches = atomic.LoadUint64(&s.ingestBatches)
	out.Ingest.Written = atomic.LoadUint64(&s.ingestWrittenRegs)
	out.Ingest.Rejected = atomic.LoadUint64(&s.ingestRejected)
	return out
}
