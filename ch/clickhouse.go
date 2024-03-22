package ch

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/danthegoodman1/RaftHouse/gologger"
	"github.com/danthegoodman1/RaftHouse/utils"
	"time"
)

var (
	RWConn driver.Conn
	RConn  driver.Conn

	logger = gologger.NewLogger()
)

const (
	RaftIndexID = "raft_idx"
)

func init() {
	s := time.Now()
	var err error
	RWConn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{utils.CH_WRITE_DSN},
		// Debug:           true,
		DialTimeout:     time.Second * 5,
		MaxOpenConns:    10,
		MaxIdleConns:    0,
		ConnMaxLifetime: time.Hour,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		TLS: &tls.Config{},
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to clickhouse write DSN")
	}

	RConn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{utils.CH_READ_DSN},
		// Debug:           true,
		DialTimeout:     time.Second * 5,
		MaxOpenConns:    10,
		MaxIdleConns:    0,
		ConnMaxLifetime: time.Hour,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		TLS: &tls.Config{},
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to clickhouse write DSN")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ctx = clickhouse.Context(ctx, clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}))

	if err := RWConn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			logger.Fatal().Msg(fmt.Sprintf("error pinging RW Clickhouse: %d %s \n%s", exception.Code, exception.Message, exception.StackTrace))
		}
		logger.Fatal().Err(err).Msg("error pinging RW clickhouse")
	}

	if err := RConn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			logger.Fatal().Msg(fmt.Sprintf("error pinging R Clickhouse: %d %s \n%s", exception.Code, exception.Message, exception.StackTrace))
		}
		logger.Fatal().Err(err).Msg("error pinging R clickhouse")
	}

	logger.Debug().Msgf("connected to clickhouse in %s", time.Since(s))
}
