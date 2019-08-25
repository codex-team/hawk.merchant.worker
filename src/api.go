package main

import (
	"github.com/codex-team/tinkoff.api.golang"
	"log"
)

func confirm(request tinkoff.ConfirmRequest) {
	client := tinkoff.NewClient("1545087497304DEMO", "9s2jgyxjhvuwxcba")
	result, err := client.Confirm(&request)
	if err != nil {
		log.Printf("Confirmation error: %s", err)
	}

	log.Printf("%s", result)
}
