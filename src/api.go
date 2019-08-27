package main

import (
	"github.com/codex-team/tinkoff.api.golang"
	"log"
)

func confirm(request tinkoff.ConfirmRequest) error {
	client := tinkoff.NewClient(tinkoffTerminalKey, tinkoffSecretKey)
	_, err := client.Confirm(&request)
	if err != nil {
		log.Printf("Confirmation error: %s", err)
	}

	return err
}
