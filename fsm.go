package raftwrapper

import (
	"fmt"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/gorilla/mux"
)

type FSM interface {
	Apply(entry raftpb.Entry)
	Snapshot() (raftpb.Snapshot, error)
	Restore(snap raftpb.Snapshot) error
	RegisterAPI(router *mux.Router)
}
