# Terraform Provider Integration Tests

This directory contains an integration test framework for the Verity Terraform provider. The framework automatically discovers `.tf` files in the directory where it's run and tests them using test cases defined in separate `.tf` files.

## Requirements

Before running the tests, install the required Python packages:

```bash
cd test_integration
pip install -r requirements.txt
```

This will install:
- **pytest** — Test framework
- **python-hcl2** — HCL parser for partial resource definitions and merging

## Overview

The test framework supports two kinds of resources:

### Normal Resources (add.tf + modify.tf)

Resources that support create and delete operations. The test lifecycle is:
1. **Phase 1 — Create**: Inject new resource blocks from `add.tf` into the `.tf` file, run `terraform apply`
2. **Phase 2 — Modify**: Replace injected blocks with a merged version of `add.tf` + `modify.tf`, run `terraform apply`
3. **Phase 3 — Cleanup**: Remove injected blocks, run `terraform apply` to delete the test resources

### Update-Only Resources (modify.tf only)

Resources that already exist on the server and cannot be created or deleted (e.g., `sites`, `sfpbreakouts`). The test lifecycle is:
1. **Phase 1 — Modify**: Merge `modify.tf` overrides into the existing resource block in the `.tf` file in-place, run `terraform apply`
2. **Phase 2 — Revert**: Restore the original `.tf` file from backup, run `terraform apply` to revert the server to its original state

A resource is automatically detected as **update-only** when its `test_cases/` directory contains `modify.tf` but **no** `add.tf`.

## Provider Mode

The Verity provider supports two modes: **datacenter** and **campus**. Different modes may expose different resources or different schemas for the same resource. Test cases are organized into separate folders per mode.

The `--mode` flag is **mandatory** for all test runs.

## Directory Structure

```
test_integration/
├── test_terraform_resources.py     # Main test runner
├── requirements.txt                # Python dependencies
├── README.md                       # This file
└── test_cases/
    ├── datacenter/                 # Test cases for datacenter mode
    │   ├── gateways/
    │   │   ├── add.tf
    │   │   └── modify.tf
    │   ├── tenants/
    │   │   ├── add.tf
    │   │   └── modify.tf
    │   ├── sites/                  # Update-only (no add.tf)
    │   │   └── modify.tf
    │   └── ...
    └── campus/                     # Test cases for campus mode
        ├── tenants/
        │   ├── add.tf
        │   └── modify.tf
        └── ...
```

## Usage

### Prerequisites

Install Python dependencies:

```bash
cd test_integration
pip install -r requirements.txt
```

Using a virtual environment (recommended):
```bash
python3 -m venv test_integration/venv
source test_integration/venv/bin/activate
pip install -r test_integration/requirements.txt
```

Make sure your Terraform environment is already initialized (the script does **not** run `terraform init`).

### Running Tests

Run the script from the directory containing your `.tf` files (e.g., `examples/`). The `--mode` flag is **required**:

```bash
cd examples

# Run the full batch test (datacenter mode)
pytest ../test_integration/test_terraform_resources.py -v -s --mode datacenter

# Run a single resource test (datacenter mode)
pytest ../test_integration/test_terraform_resources.py -v -s --mode datacenter -k badges

# Run campus mode tests
pytest ../test_integration/test_terraform_resources.py -v -s --mode campus

# Run a single resource test (campus mode)
pytest ../test_integration/test_terraform_resources.py -v -s --mode campus -k tenants
```

### Test Modes

**Batch mode (default)** — `test_all_resources_lifecycle`

Runs all discovered resources together. Each phase (create, modify, cleanup) issues a single `terraform apply` covering all resources at once. This verifies that resources work correctly in combination and respects operation ordering.

**Single-resource mode** — `test_single_resource[name]` (via `-k`)

Runs the full lifecycle for one resource in isolation. Useful for debugging a specific resource or verifying a new test case. Single-resource tests are automatically deselected unless you pass `-k`.

### Available pytest flags

| Flag | Description |
|------|-------------|
| `--mode <mode>` | **Required.** Provider mode: `datacenter` or `campus` |
| `-v` | Verbose output (show test details) |
| `-s` | Show print statements and terraform output (**highly recommended**) |
| `-x` | Stop on first failure |
| `-k <keyword>` | Run single-resource tests matching the keyword |

## Creating Test Cases

### Normal Resource (create + modify + delete)

1. Create a directory matching the `.tf` file name (without extension) under the appropriate mode:
   ```bash
   mkdir -p test_integration/test_cases/datacenter/gateways
   ```

2. Create `add.tf` with the full resource definitions to create:
   ```hcl
   resource "verity_gateway" "gateway_test_1" {
     name                      = "gateway_test_1"
     depends_on                = [verity_operation_stage.gateway_stage]
     advertisement_interval    = 30
     bfd_transmission_interval = 300
     keepalive_timer           = 60
     neighbor_ip_address       = "8.8.8.8"
     tenant                    = "Visualizator"
     tenant_ref_type_          = "tenant"
     object_properties {
       group = ""
     }
   }
   ```

3. Create `modify.tf` with **only the fields you want to change**:
   ```hcl
   resource "verity_gateway" "gateway_test_1" {
     bfd_transmission_interval = 311
     keepalive_timer           = 66
   }
   ```

   The framework merges `modify.tf` on top of `add.tf` — all unspecified fields are preserved automatically.

### Update-Only Resource (modify + revert)

For resources that already exist and cannot be created or deleted:

1. Create a directory with **only** `modify.tf` (no `add.tf`) under the appropriate mode:
   ```bash
   mkdir -p test_integration/test_cases/datacenter/sites
   ```

2. Create `modify.tf` referencing the **existing** resource by its type and name as they appear in the `.tf` file:
   ```hcl
   resource "verity_site" "some_existing_site" {
     description = "Modified by test framework"
   }
   ```

   The framework reads the current resource block from the `.tf` file, merges the overrides on top, and writes the result back in-place.

## Merge Behavior

The `modify.tf` files support partial definitions. The merge follows these rules:

| Scenario | Behavior |
|----------|----------|
| **Top-level attributes** | Override wins — specify only changed fields |
| **Nested blocks (dict + dict)** | Recursive deep merge |
| **Indexed nested blocks** (list of dicts with `index` key) | Merge by matching `index` value — unmentioned indexes preserved |
| **Non-indexed nested blocks** (list of dicts without `index`) | Full replacement — specify all entries |
| **Simple lists** (strings, numbers) | Full replacement |

### Example: Indexed block merge

If the base resource has `route_tenants` blocks with indexes 1, 2, 3, and `modify.tf` specifies only index 2 with a changed field, indexes 1 and 3 are preserved unchanged and index 2 is deep-merged.

### Example: Non-indexed block

For blocks like `object_properties` (parsed as a list of one dict without `index`), the entire block is replaced. Specify all fields from the original in `modify.tf`.

## Resource Discovery

The framework discovers resources using these rules:

1. Scans the current directory for `*.tf` files
2. Filters to only resources listed in `internal/importer/importer.go`
3. Checks `test_cases/{mode}/{name}/` for `add.tf` and/or `modify.tf` with valid resource definitions
4. Classifies each as **normal** (has `add.tf`) or **update-only** (has only `modify.tf`)
5. Skips resources with empty or comment-only test files

## Batch Test Lifecycle

### Pre-flight Check
- Runs `terraform plan` to verify a clean starting state
- Fails immediately if there are pending changes

### Phase 1 — Batch Create (normal resources only)
1. Backs up all `.tf` files
2. Appends resources from `add.tf` into each `.tf` file (between `# === TEST RESOURCES START/END ===` markers)
3. Runs a single `terraform apply` for all new resources

### Phase 2 — Batch Modify (all resources)
1. **Normal resources**: Removes injected blocks, merges `modify.tf` with `add.tf`, re-injects the merged result
2. **Update-only resources**: Merges `modify.tf` overrides into the existing resource block in the `.tf` file in-place
3. Runs a single `terraform apply` for all modifications

### Phase 3 — Batch Cleanup
1. **Normal resources**: Removes injected blocks (deletes the test resources)
2. **Update-only resources**: Restores original `.tf` file from backup (reverts to original state)
3. Runs a single `terraform apply` to delete/revert everything
4. Restores all `.tf` files from backup

On failure at any point, all `.tf` files are restored from backup automatically.

## API Request Logging

Every test run automatically captures and logs all **PUT**, **PATCH**, and **DELETE** HTTP requests made by the Terraform provider during `terraform apply`. This lets you verify request ordering, payload correctness, and response status without manually sifting through terraform debug output.

### How It Works

1. **Debug capture**: Each `terraform apply` invocation runs with `TF_LOG=DEBUG` and `TF_LOG_PATH` pointed at a temporary file. This captures the full provider-level HTTP traffic.

2. **Parsing**: After the apply completes, the framework parses the debug log with regex to extract HTTP exchanges. For each PUT/PATCH/DELETE request it captures:
   - *Timestamp* — when the request was sent
   - *Method & URL* — e.g. `PUT /tenants/test_tenant_1`
   - *Request body* — the JSON payload sent to the server
   - *Response status* — e.g. `200 OK`
   - *Response body* — the JSON returned by the server

   GET requests are excluded to keep the log focused on state-changing operations.

3. **Logging**: Extracted exchanges are appended to a timestamped log file grouped by phase label (e.g. "Phase 1: Add [tenants]", "Phase 2: Modify", "Phase 3: Cleanup"). The temporary debug file is deleted after parsing.

### Log File

Each test run creates a unique log file in the working directory (e.g. `examples/`):

```
test_logs_20260227_171624.log
test_logs_20260227_183045.log
```

The filename includes a timestamp (`YYYYMMDD_HHMMSS`) so that logs from different runs are never overwritten.

### Log Format

```
================================================================================
API Request Log: Single Resource: tenants
================================================================================

────────────────────────────────────────────────────────────────────────────────
Phase 1: Add [tenants]
────────────────────────────────────────────────────────────────────────────────

[17:21:28] PUT /api/tenants
  → {"tenant":{"tenant_test_script1":{"default_originate":false,"enable":true,"name":"tenant_test_script1",...}}}
  ← 201 Created
  ← {"status": "succeeded", "message": "Tenant 'tenant_test_script1' successfully created"}

────────────────────────────────────────────────────────────────────────────────
Phase 2: Modify [tenants]
────────────────────────────────────────────────────────────────────────────────

[17:24:33] PATCH /api/tenants
  → {"tenant":{"tenant_test_script1":{"default_originate":true,"enable":false,"object_properties":{"group":"test"},...}}}
  ← 200 OK
  ← {"status": "succeeded", "message": "Tenant 'tenant_test_script1' successfully updated"}

────────────────────────────────────────────────────────────────────────────────
Phase 3: Cleanup [tenants]
────────────────────────────────────────────────────────────────────────────────

[17:27:29] DELETE /api/tenants?tenant_name=tenant_test_script1
  ← 200 OK
  ← {"status": "succeeded", "message": "Tenants tenant_test_script1 successfully deleted"}
```

## Valid Resource Names

The script only tests resources defined in `internal/importer/importer.go`. If a `.tf` file doesn't match a name in that list, it is skipped.

