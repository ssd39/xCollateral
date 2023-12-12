package main

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	log.Level(zerolog.Disabled)
	os.Exit(m.Run())
}
