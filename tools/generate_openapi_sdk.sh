#!/usr/bin/env bash
set -euo pipefail

readonly task_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly task_inputs="$task_root/specs/openapi/6.6"
readonly task_generator_image="openapitools/openapi-generator-cli@sha256:2ab0a9680222de65dc9d3baf861aa02b99e1b80c211d8221ebf3ae8f8a102524"
readonly task_target="$task_root/openapi"

case "${1:---check}" in
  --check|--write) task_mode="${1:---check}" ;;
  *) echo "usage: $0 [--check|--write]" >&2; exit 2 ;;
esac

go run "$task_root/tools/specgen" verify --input-dir "$task_inputs"

task_temp="$(mktemp -d)"
trap 'rm -rf "$task_temp"' EXIT

python3 "$task_root/tools/process_swagger.py" \
  "$task_inputs/datacenter.json" \
  "$task_inputs/campus.json" \
  --output "$task_temp/openapi.json"

docker run --rm \
  --user "$(id -u):$(id -g)" \
  --volume "$task_temp:/work" \
  "$task_generator_image" generate \
  --input-spec /work/openapi.json \
  --generator-name go \
  --output /work/sdk \
  --global-property apiTests=false,modelTests=false \
  --additional-properties packageName=openapi

if [[ "$task_mode" == "--check" ]]; then
	if ! diff -qr "$task_target" "$task_temp/sdk"; then
    echo "OpenAPI SDK drift detected. Review and apply with: tools/generate_openapi_sdk.sh --write" >&2
    exit 1
  fi
  exit 0
fi

rsync -a --delete "$task_temp/sdk/" "$task_target/"
