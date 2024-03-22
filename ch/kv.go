package ch

import (
	"github.com/danthegoodman1/RaftHouse/utils"
	"github.com/dgraph-io/badger/v4"
)

var (
	KV *badger.DB
)

func init() {
	logger.Debug().Msgf("opening DB path at %s", utils.DB_PATH)
	var err error
	KV, err = badger.Open(badger.DefaultOptions(utils.DB_PATH))
	if err != nil {
		logger.Fatal().Err(err).Msg("error in badger.Open")
	}
}
