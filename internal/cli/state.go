package cli

import (
	"github.com/abahnj/rssagg/internal/config"
)

// State holds application state including configuration
type State struct {
	Config *config.Config
}