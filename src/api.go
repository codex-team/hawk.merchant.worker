package main

import (
	"github.com/codex-team/tinkoff.api.golang"
	"log"
)

func confirm(request tinkoff.ConfirmRequest) {
	client := tinkoff.NewClient(tinkoffTerminalKey, tinkoffSecretKey)
	result, err := client.Confirm(&request)
	if err != nil {
		log.Printf("Confirmation error: %s", err)
	}

	log.Printf("%s", result)
}
