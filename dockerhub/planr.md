# Planr

Grundrisse und Häuser maßstabsgetreu planen, im Browser. Wohnungen oder ganze Häuser zeichnen, Wände, Türen und Fenster setzen, möblieren und Flächen berechnen. Konten, Freigaben und optional Anmeldung über OpenID Connect (SSO). Alles steckt in einem Container, die Daten liegen in `/data`. Eine Mac-App gibt es bei den Releases auf GitHub.

Quellcode und Anleitung: [https://github.com/Dschonas04/Planr](https://github.com/Dschonas04/Planr)

## Start

```bash
docker run -d --name planr -p 8090:8090 -v planr:/data dschohnas/planr:2
```

Danach läuft Planr unter `http://<host>:8090`.

## Tags

`latest`, die genaue Version (`2.1.0`) sowie `2.1` und `2`. Jedes Image gibt es für `linux/amd64` und `linux/arm64`. Dieselben Images liegen auch unter `ghcr.io/dschonas04/planr`.

## Lizenz

Business Source License 1.1. Lesen, ändern, selbst betreiben und in der eigenen Organisation nutzen ist erlaubt, nicht aber der Betrieb oder Verkauf als kommerzielles Angebot für Dritte. Am 10. August 2030 geht diese Version in Apache 2.0 über.
