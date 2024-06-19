package main

import "testing"

func TestFormatMeasurements(t *testing.T) {
	measurements := make(map[string]*Stats)
	measurements["Abha"] = &Stats{Min: -23.0, Max: 18.0, Mean: 59.2}
	measurements["Yerevan"] = &Stats{Min: -37.2, Max: 12.4, Mean: 59.0}
	measurements["Pontianak"] = &Stats{Min: -23.3, Max: 27.7, Mean: 76.4}
	measurements["Bosaso"] = &Stats{Min: -19.0, Max: 30.0, Mean: 78.5}

	results := formatMeasurements(measurements)
	expected := "{Abha=-23.0/59.2/18.0, Bosaso=-19.0/78.5/30.0, Pontianak=-23.3/76.4/27.7, Yerevan=-37.2/59.0/12.4}"

	if results != expected {
		t.Errorf("Result does not match expected. result=%s, expected=%s", results, expected)
	}
}
