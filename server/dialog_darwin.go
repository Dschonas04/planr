//go:build darwin

package main

// Dateidialoge ueber AppleScript. osascript gehoert zu jedem macOS, braucht
// kein cgo und zeigt die gewohnten Sichern-/Oeffnen-Fenster. Der
// Dateivorschlag geht als Argument hinein, nicht in den Skripttext -- ein
// Name mit Anfuehrungszeichen kann so nichts ausfuehren.

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type appleScriptDialoge struct{}

var systemDialoge dateiDialoge = appleScriptDialoge{}

func osascript(zeilen []string, args ...string) (string, error) {
	cmdArgs := []string{}
	for _, z := range zeilen {
		cmdArgs = append(cmdArgs, "-e", z)
	}
	cmdArgs = append(cmdArgs, args...)
	out, err := exec.Command("osascript", cmdArgs...).Output()
	if err != nil {
		var ee *exec.ExitError
		// -128 = vom Nutzer abgebrochen
		if errors.As(err, &ee) && strings.Contains(string(ee.Stderr), "-128") {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Wie unter Linux lassen sich die Dialoge fuer Tests per Umgebung ersetzen;
// auf einem Mac ohne Bildschirm (CI) wartete osascript sonst ewig.
func (appleScriptDialoge) Speichern(vorschlag string) (string, error) {
	if verz := os.Getenv("PLANR_DIALOG_VERZ"); verz != "" {
		return filepath.Join(verz, vorschlag), nil
	}
	return osascript([]string{
		"on run argv",
		"activate",
		`set f to choose file name with prompt "Sichern unter" default name (item 1 of argv) default location (path to documents folder)`,
		"return POSIX path of f",
		"end run",
	}, vorschlag)
}

func (appleScriptDialoge) Oeffnen() (string, error) {
	if pfad := os.Getenv("PLANR_DIALOG_OEFFNEN"); pfad != "" {
		return pfad, nil
	}
	return osascript([]string{
		"activate",
		`set f to choose file with prompt "Planr-Datei öffnen (.planr oder .json)" default location (path to documents folder)`,
		"return POSIX path of f",
	})
}
