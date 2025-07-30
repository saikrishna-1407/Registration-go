package handlers

import (
	"finalreg/config"
	"finalreg/internal/providers"
)

type Service struct {
	ServiceName string
	Config      *config.Config
	Db          providers.RepoStore
}
