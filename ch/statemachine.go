package ch

import (
	"encoding/binary"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/dgraph-io/badger/v4"
	"github.com/lni/dragonboat/v3/statemachine"
	"io"
)

type CHStateMachine struct {
	conn driver.Conn
}

func NewCHStateMachine(conn driver.Conn) statemachine.IOnDiskStateMachine {
	machine := &CHStateMachine{
		conn: conn,
	}
	return machine
}

func (c *CHStateMachine) Open(stopc <-chan struct{}) (uint64, error) {
	var idx uint64
	err := KV.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(RaftIndexID))
		if err != nil {
			return fmt.Errorf("error in txn.Get: %w", err)
		}
		return item.Value(func(val []byte) error {
			binary.LittleEndian.PutUint64(val, idx)
			return nil
		})
	})
	if err != nil {
		return 0, fmt.Errorf("error in getting raft index: %w", err)
	}

	return idx, nil
}

func (c *CHStateMachine) Update(entries []statemachine.Entry) ([]statemachine.Entry, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CHStateMachine) Lookup(i interface{}) (interface{}, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CHStateMachine) Sync() error {
	// No-op
	return nil
}

func (c *CHStateMachine) PrepareSnapshot() (interface{}, error) {
	// TODO create zip
	panic("implement me")
}

func (c *CHStateMachine) SaveSnapshot(i interface{}, writer io.Writer, i2 <-chan struct{}) error {
	// TODO read zip
	panic("implement me")
}

func (c *CHStateMachine) RecoverFromSnapshot(reader io.Reader, i <-chan struct{}) error {
	// TODO download zip and recover
	panic("implement me")
}

func (c *CHStateMachine) Close() error {
	// No-op, we close connections on node shutdown
	return nil
}
