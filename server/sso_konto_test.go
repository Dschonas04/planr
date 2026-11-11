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
	eins, neu, err := k.SSOKonto("Anna@Example.org", "Anna", "iss|a", false, true)
	if err != nil || !neu || eins.Rolle != RolleAdmin || eins.Email != "anna@example.org" {
		t.Fatalf("erstes Konto: %+v %v %v", eins, neu, err)
	}
	wieder, neu, err := k.SSOKonto("anna@example.org", "", "iss|a", false, true)
	if err != nil || neu || wieder.ID != eins.ID {
		t.Fatalf("Konto nicht wiedergefunden: %+v %v %v", wieder, neu, err)
	}
	if _, _, err := k.SSOKonto("dora@example.org", "", "iss|d", false, false); err != ErrSSOKeinKonto {
		t.Fatalf("Konto trotz geschlossener Registrierung angelegt: %v", err)
	}
	zwei, _, _ := k.SSOKonto("bert@example.org", "Bert", "iss|b", false, true)
	if zwei.Rolle != RolleNutzer {
		t.Fatal("zweites Konto ist Administrator")
	}
	gesperrt := true
	if _, err := k.Aendern(zwei.ID, nil, &gesperrt); err != nil {
		t.Fatal(err)
	}
	if _, _, err := k.SSOKonto("anna@example.org", "", "iss|fremd", false, true); err != ErrSSOFremd {
		t.Fatalf("fremdes Anbieter-Konto angenommen: %v", err)
	}
	// Vorab mit Passwort angelegt, Registrierung offen: Passwort und Sitzungen weg.
	vorab, _ := k.Anlegen("carl@example.org", "", "passwort-123456", RolleNutzer)
	token := k.SitzungAnlegen(vorab.ID)
	if _, _, err := k.SSOKonto("carl@example.org", "", "iss|c", true, true); err != nil {
		t.Fatal(err)
	}
	if _, err := k.Pruefen("carl@example.org", "passwort-123456"); err == nil {
		t.Fatal("vorab gesetztes Passwort gilt noch")
	}
	if _, ok := k.SitzungKonto(token); ok {
		t.Fatal("alte Sitzung gilt noch")
	}
	// Früher selbst registriert, Registrierung inzwischen zu: trotzdem zurückgesetzt.
	selbst, _ := k.Anlegen("eva@example.org", "", "passwort-123456", RolleNutzer)
	k.SelbstMarkieren(selbst.ID)
	if _, _, err := k.SSOKonto("eva@example.org", "", "iss|e", false, false); err != nil {
		t.Fatal(err)
	}
	if _, err := k.Pruefen("eva@example.org", "passwort-123456"); err == nil {
		t.Fatal("Passwort eines selbst registrierten Kontos gilt noch")
	}
	if _, _, err := k.SSOKonto("bert@example.org", "", "iss|b", false, true); err != ErrKontoGesperrt {
		t.Fatalf("gesperrtes Konto angemeldet: %v", err)
	}
}
