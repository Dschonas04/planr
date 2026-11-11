package main

import (
	"testing"
	"time"
)

func TestSSOKonto(t *testing.T) {
	k, err := NeueKonten(t.TempDir(), time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	eins, neu, err := k.SSOKonto("Anna@Example.org", "Anna")
	if err != nil || !neu || eins.Rolle != RolleAdmin || eins.Email != "anna@example.org" {
		t.Fatalf("erstes Konto: %+v %v %v", eins, neu, err)
	}
	wieder, neu, err := k.SSOKonto("anna@example.org", "")
	if err != nil || neu || wieder.ID != eins.ID {
		t.Fatalf("Konto nicht wiedergefunden: %+v %v %v", wieder, neu, err)
	}
	zwei, _, _ := k.SSOKonto("bert@example.org", "Bert")
	if zwei.Rolle != RolleNutzer {
		t.Fatal("zweites Konto ist Administrator")
	}
	gesperrt := true
	if _, err := k.Aendern(zwei.ID, nil, &gesperrt); err != nil {
		t.Fatal(err)
	}
	if _, _, err := k.SSOKonto("bert@example.org", ""); err != ErrKontoGesperrt {
		t.Fatalf("gesperrtes Konto angemeldet: %v", err)
	}
}
