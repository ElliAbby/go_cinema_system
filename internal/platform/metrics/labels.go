package metrics

import (
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func labeledRegisterer(serviceName string) prometheus.Registerer {
	return prometheus.WrapRegistererWith(prometheus.Labels{
		"app":      serviceName,
		"env":      envOrDefault("APP_ENV", "ENVIRONMENT", "local"),
		"instance": envOrDefault("INSTANCE", "HOSTNAME", "local"),
	}, prometheus.DefaultRegisterer)
}

func envOrDefault(keys ...string) string {
	defaultValue := ""
	if len(keys) > 0 {
		defaultValue = keys[len(keys)-1]
		keys = keys[:len(keys)-1]
	}

	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			trimmed := strings.TrimSpace(value)
			if trimmed != "" {
				return trimmed
			}
		}
	}

	return defaultValue
}