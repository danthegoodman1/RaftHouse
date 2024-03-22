package ch

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/lni/dragonboat/v3/statemachine"
	"io"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	row := c.conn.QueryRow(ctx, "select coalesce(max(idx), 0) from raft_index where id = 0")
	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("error in QueryRow: %w", err)
	}

	var idx uint64
	if err := row.Scan(&idx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// New DB
			return 0, nil
		}
		return 0, fmt.Errorf("error in row.Scan: %w", err)
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
