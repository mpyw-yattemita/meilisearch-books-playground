// Package config reads settings shared by the commands from environment variables.
package config

import (
	"os"

	"github.com/meilisearch/meilisearch-go"
)

// Defaults match compose.yaml.
const (
	defaultMeiliHost = "http://localhost:7700"
	defaultMeiliKey  = "local-dev-master-key-change-me"
)

// NewMeilisearchClient returns a client for MEILI_HOST using MEILI_MASTER_KEY.
func NewMeilisearchClient() meilisearch.ServiceManager {
	return meilisearch.New(env("MEILI_HOST", defaultMeiliHost), meilisearch.WithAPIKey(env("MEILI_MASTER_KEY", defaultMeiliKey)))
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
