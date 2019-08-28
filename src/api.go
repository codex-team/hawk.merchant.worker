package main

import (
	"github.com/codex-team/tinkoff.api.golang"
	"log"
)

func confirm(request tinkoff.ConfirmRequest) (tinkoff.ConfirmResponse, error) {
	client := tinkoff.NewClient(tinkoffTerminalKey, tinkoffSecretKey)
	resp, err := client.Confirm(&request)
	if err != nil {
		log.Printf("Confirmation error: %s", err)
	}

	return resp, err
}
