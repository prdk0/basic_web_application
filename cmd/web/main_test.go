package main

import "testing"

func TestRun(t *testing.T) {
	db, err := run()
	if err != nil {
		t.Error("failed to run run()")
	}
	defer db.SQL.Close()
}
