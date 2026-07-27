#!/bin/bash
set -euo pipefail

ARCH="x86"
GOARCH="386"
OUTPUT="dist/ThunderUpdater-x86.exe"
SYSO="cmd/thunderupdater/resource_windows_386.syso"

echo "=== Build Windows $ARCH (32-bit) ==="

mkdir -p dist

if ! command -v goversioninfo &>/dev/null; then
	echo "Instalando goversioninfo..."
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
	export PATH="$PATH:$(go env GOPATH)/bin"
else
	echo "goversioninfo encontrado."
fi

echo "Gerando resource.syso para 386..."
goversioninfo -64=false -icon=assets/icon.ico -o "$SYSO" cmd/thunderupdater/versioninfo.json

CGO_ENABLED=1 GOOS=windows GOARCH="$GOARCH" \
  go build -ldflags="-s -w" -o "$OUTPUT" ./cmd/thunderupdater

rm -f "$SYSO"

SIZE=$(stat -c%s "$OUTPUT" 2>/dev/null || stat -f%z "$OUTPUT" 2>/dev/null)
SIZE_MB=$(echo "scale=2; $SIZE / 1048576" | bc 2>/dev/null || echo "$SIZE bytes")

echo ""
echo "--- Build concluído ---"
echo "Arquitetura: $ARCH ($GOARCH)"
echo "Tamanho:     $SIZE_MB"
echo "Arquivo:     $OUTPUT"
