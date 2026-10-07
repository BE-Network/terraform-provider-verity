#!/usr/bin/env bash
set -euo pipefail

readonly task_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
readonly task_generator_image="openapitools/openapi-generator-cli@sha256:2ab0a9680222de65dc9d3baf861aa02b99e1b80c211d8221ebf3ae8f8a102524"
task_mode=--check
task_mode_set=false
task_sdk_only=false
for task_arg in "$@"; do
  case "$task_arg" in
    --check|--write)
      if "$task_mode_set"; then
        echo "choose either --check or --write" >&2
        exit 2
      fi
      task_mode="$task_arg"
      task_mode_set=true
      ;;
    --sdk-only) task_sdk_only=true ;;
    --help)
      echo "usage: $0 [--check|--write] [--sdk-only]"
      exit 0
      ;;
    *) echo "usage: $0 [--check|--write] [--sdk-only]" >&2; exit 2 ;;
  esac
done

cd "$task_root"
task_temp="$(mktemp -d)"
trap 'rm -rf "$task_temp"' EXIT
go build -o "$task_temp/specgen" ./tools/specgen
go build -o "$task_temp/sdkgen" ./tools/sdkgen
task_version="$("$task_temp/specgen" version --overrides specs/overrides.yaml)"
readonly task_version
readonly task_inputs="$task_root/specs/openapi/$task_version"
"$task_temp/specgen" verify --input-dir "$task_inputs" --api-version "$task_version"

if [[ "$task_mode" == "--write" ]] && ! "$task_sdk_only"; then
  "$task_temp/specgen" extract --input-dir "$task_inputs" --output "$task_temp/manifest.json"
  mkdir -p "$task_temp/embed"
  "$task_temp/specgen" registry --input-dir "$task_inputs" --overrides specs/overrides.yaml \
    --output "$task_temp/registry.json" --embed-output "$task_temp/embed/registry.json"
fi

python3 tools/process_swagger.py "$task_inputs/datacenter.json" "$task_inputs/campus.json" \
  --output "$task_temp/openapi.json"
docker run --rm --user "$(id -u):$(id -g)" --volume "$task_temp:/work" \
  "$task_generator_image" generate --input-spec /work/openapi.json --generator-name go \
  --output /work/sdk --global-property apiTests=false,modelTests=false \
  --additional-properties packageName=openapi
"$task_temp/sdkgen" --sdk-dir "$task_temp/sdk"

if [[ "$task_mode" == "--check" ]]; then
  if ! diff -qr openapi "$task_temp/sdk"; then
    echo "OpenAPI SDK drift detected. Review and regenerate with: tools/generate_provider.sh --write" >&2
    exit 1
  fi
  if "$task_sdk_only"; then
    exit 0
  fi
  "$task_temp/specgen" extract --input-dir "$task_inputs" --output specs/generated_manifest.json --check
  "$task_temp/specgen" registry --input-dir "$task_inputs" --overrides specs/overrides.yaml \
    --output specs/generated_registry.json --embed-output internal/registry/registry.json --check
  "$task_temp/specgen" adapters --registry specs/generated_registry.json --openapi-dir openapi \
    --output internal/transport/generated_adapters.go --check
  "$task_temp/specgen" bulk --registry specs/generated_registry.json --openapi-dir openapi \
    --output internal/bulkops/generated_registry.go --check
  "$task_temp/specgen" docs --registry specs/generated_registry.json --output-dir docs/resources --check
  exit 0
fi

if ! "$task_sdk_only"; then
  "$task_temp/specgen" adapters --registry "$task_temp/registry.json" --openapi-dir "$task_temp/sdk" \
    --output "$task_temp/generated_adapters.go"
  "$task_temp/specgen" bulk --registry "$task_temp/registry.json" --openapi-dir "$task_temp/sdk" \
    --output "$task_temp/generated_bulk_registry.go"
  mkdir -p "$task_temp/docs"
  cp docs/resources/verity_operation_stage.md "$task_temp/docs/"
  "$task_temp/specgen" docs --registry "$task_temp/registry.json" --output-dir "$task_temp/docs"
fi

rsync -a --delete "$task_temp/sdk/" openapi/
if ! "$task_sdk_only"; then
  cp "$task_temp/manifest.json" specs/generated_manifest.json
  cp "$task_temp/registry.json" specs/generated_registry.json
  cp "$task_temp/embed/registry.json" internal/registry/registry.json
  cp "$task_temp/generated_adapters.go" internal/transport/generated_adapters.go
  cp "$task_temp/generated_bulk_registry.go" internal/bulkops/generated_registry.go
  rsync -a --delete "$task_temp/docs/" docs/resources/
fi
