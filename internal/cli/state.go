package cli

import (
	"github.com/abahnj/rssagg/internal/config"
	"github.com/abahnj/rssagg/internal/database"
)

// State holds application state including configuration
type State struct {
	Config *config.Config
	Db 	*database.Queries
}