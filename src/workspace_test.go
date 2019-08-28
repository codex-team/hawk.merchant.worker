package main

import (
	"testing"
)

func TestWorkspaceInc(t *testing.T) {
	loadEnv()
	database := connectMongo()
	err := updateWorkspaceBalance(database, "5d655e57a62cda909a438b8d", 10)
	if err != nil {
		t.Error(err)
	}
}
