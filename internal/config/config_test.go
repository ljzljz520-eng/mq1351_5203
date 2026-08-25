package config

import "testing"

func TestConfigDefaults(t *testing.T) {
	c := LoadFrom(func(string) string { return "" })
	if c.Address != ":8080" || c.Database != "qin-culture.db" || c.MaxBody <= 0 {
		t.Fatalf("unexpected defaults: %#v", c)
	}
}
