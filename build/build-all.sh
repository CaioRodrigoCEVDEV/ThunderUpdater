#!/bin/bash
set -euo pipefail

echo "========================================"
echo "  ThunderUpdaterGO - Build All"
echo "========================================"
echo ""

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

"$SCRIPT_DIR/build-x86.sh"
echo ""
"$SCRIPT_DIR/build-x64.sh"

echo ""
echo "========================================"
echo "  Todos os builds concluídos"
echo "========================================"
