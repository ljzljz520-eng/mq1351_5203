package main

import (
	"testing"

	"qin-culture-site/internal/config"
)

func TestRunRejectsInvalidConfig(t *testing.T) {
	if err := run(config.Config{}); err == nil {
		t.Fatal("expected invalid config error")
	}
}
