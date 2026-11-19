package main

import "testing"

func TestEchteIP(t *testing.T) {
	faelle := []struct {
		direkt string
		xff    []string
		will   string
	}{
		{"10.0.2.6", []string{"203.0.113.7"}, "203.0.113.7"},
		// Besucher schreibt selbst etwas davor: zaehlt nicht
		{"10.0.2.6", []string{"1.2.3.4, 203.0.113.7, 10.0.2.6"}, "203.0.113.7"},
		{"10.0.2.6", []string{"1.2.3.4", "203.0.113.7"}, "203.0.113.7"},
		// alles im LAN: der linke
		{"127.0.0.1", []string{"10.0.0.5, 10.0.2.6"}, "10.0.0.5"},
		// direkt von aussen: Kopf wird nicht geglaubt
		{"198.51.100.9", []string{"10.0.0.5"}, "198.51.100.9"},
		{"10.0.2.6", nil, "10.0.2.6"},
		{"10.0.2.6", []string{"203.0.113.7, quatsch, 10.0.0.9"}, "10.0.0.9"},
	}
	for _, f := range faelle {
		if g := echteIP(f.direkt, f.xff); g != f.will {
			t.Errorf("%s %v: %s, will %s", f.direkt, f.xff, g, f.will)
		}
	}
}
