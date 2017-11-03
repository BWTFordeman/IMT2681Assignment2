package main

import (
	"testing"
)

func TestFindAverageOfPost(t *testing.T) {

}

func TestFindAllWebhooks(t *testing.T) {
	name := "Fordeman"
	pass := "12345"
	webhooks, err := findAllWebhooks(name, pass)
	if err != nil {
		t.Errorf("Could not find any webhooks %v", webhooks)
	}
}
