package main

import (
	"fmt"
	"testing"
)

func TestWorkspaceInc(t *testing.T) {
	loadEnv()
	database := connectMongo()
	err := updateWorkspaceBalance(database, "5d655e57a62cda909a438b8d", 10)
	fmt.Printf("%s", err)
	t.Error()
}
