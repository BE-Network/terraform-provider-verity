#!/usr/bin/env bash
set -euo pipefail

readonly task_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
exec "$task_root/tools/generate_provider.sh" --sdk-only "$@"
