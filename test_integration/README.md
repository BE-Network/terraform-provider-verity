# Terraform Provider Integration Tests

This directory contains an integration test framework for the Verity Terraform provider. The framework automatically discovers `.tf` files in the directory where it's run and tests them using test cases defined in separate `.tf` files.

## ⚠️ Requirements

Before running the tests, you **must** install the required Python packages:

```bash
cd test_integration
pip install -r requirements.txt
```

This will install:
- **pytest** - Test framework
- **python-hcl2** - HCL parser for partial resource definitions

## Overview

The test framework:
1. Discovers all `.tf` files in the **current directory** (where you run the script from)
2. Filters them to only test resources defined in `internal/importer/importer.go`
3. Finds corresponding test cases in `test_cases/` subdirectories
4. **Phase 1 - Batch Create**: Injects test resources from all `add.tf` files into their corresponding `.tf` files, runs a single `terraform apply` for all resources
5. **Phase 2 - Batch Modify**: Removes injected resources, injects resources from all `modify.tf` files (supports **partial definitions**!), runs a single `terraform apply` for all modifications
6. **Phase 3 - Batch Cleanup**: Removes all injected test resources, runs a single `terraform apply` to clean up, and restores all original files

## Directory Structure

```
test_integration/
├── test_terraform_resources.py     # Main test runner
├── test_cases/                     # Test case definitions
│   ├── gateways/                   # Test cases for gateways.tf
│   │   ├── add.tf                  # Resources to add
│   │   └── modify.tf               # Modified versions of resources
│   ├── gatewayprofiles/            # Test cases for gatewayprofiles.tf
│   │   ├── add.tf
│   │   └── modify.tf
│   ├── ethportprofiles/            # Test cases for ethportprofiles.tf
│   │   ├── add.tf
│   │   └── modify.tf
│   └── ethportsettings/            # Test cases for ethportsettings.tf
│       ├── add.tf
│       └── modify.tf
└── README.md                       # Documentation for the script
```

## Usage

### Prerequisites

**Required Python packages** - Both pytest and python-hcl2 are required:

```bash
# Navigate to test_integration directory
cd test_integration

# Install from requirements.txt
pip install -r requirements.txt
```

**Using a virtual environment (recommended):**
```bash
# Create and activate virtual environment
python3 -m venv test_integration/venv
source test_integration/venv/bin/activate  # On Linux/Mac
# or
test_integration\venv\Scripts\activate     # On Windows
```

### Running Tests

**Important**: Run the script from the directory containing your `.tf` files (e.g., `examples/` or any directory with terraform configurations).

**Note**: The script does NOT run `terraform init`. Make sure your Terraform environment is already initialized and configured (e.g., with `dev.tfrc` for local provider development).

```bash
# Navigate to your terraform files directory
cd examples

# Run all tests (recommended: verbose + show output)
pytest ../test_integration/test_terraform_resources.py -v -s
```

**Available pytest flags:**
- `-v` - Verbose output (show test details)
- `-s` - Show print statements and terraform output (highly recommended!)
- `-x` - Stop on first failure
- `-k <keyword>` - Filter tests by keyword

The script will:
- Look for `.tf` files in the current directory
- Only test files that match resource names from `internal/importer/importer.go`
- Skip any files not defined in importer.go
- **Assumes terraform is already initialized** (does not run `terraform init`)

### Creating Test Cases

To add test cases for a new resource (e.g., `gateways.tf`):

1. Create a directory in `test_cases/` matching the resource name (without `.tf`):
   ```bash
   mkdir -p test_integration/test_cases/gateways
   ```

2. Create `add.tf` with resources to add:
   ```hcl
   # test_integration/test_cases/gateways/add.tf
   resource "verity_gateway" "gateway_test_script1" {
     name = "gateway_test_script1"
     depends_on = [verity_operation_stage.gateway_stage]
     object_properties {
       group = ""
     }
     advertisement_interval = 30
     bfd_transmission_interval = 300
     keepalive_timer = 60
     neighbor_ip_address = "8.8.8.8"
     tenant = "Visualizator"
     tenant_ref_type_ = "tenant"
     # ... other required fields
   }
   ```

3. Create `modify.tf` with **only the fields you want to change**:
   ```hcl
   # test_integration/test_cases/gateways/modify.tf
   resource "verity_gateway" "gateway_test_script1" {
     bfd_transmission_interval = 311
     keepalive_timer = 66
   }
   ```

   You only need to specify the fields that change. All other fields from `add.tf` are automatically preserved and merged.
   
   **⚠️ Important**: `modify.tf` can only exist if `add.tf` exists. You cannot modify resources that weren't added in Phase 1.

4. Run the tests:
   ```bash
   pytest test_integration/test_terraform_resources.py -v -s
   ```

## Valid Resource Names

The script only tests resources defined in `internal/importer/importer.go`:
- tenants
- gateways
- gatewayprofiles
- services
- ethportprofiles
- ethportsettings
- lags
- bundles

If a `.tf` file doesn't match these names, it will be skipped.

## Test Lifecycle

The test framework processes all resources in three global phases for efficiency:

### Pre-flight Check: Verify Clean State
Before any testing begins, the framework runs a single `terraform plan` to ensure the starting state is clean:
1. Run `terraform plan` once
2. Verify "No changes" in the output
3. If there are pending changes, fail the test immediately to prevent testing on an inconsistent state

This ensures all tests start from a known good state and helps catch configuration issues early.

### Phase 1: Add All Resources
1. Backup all original `.tf` files
2. Inject resources from `test_cases/{resource}/add.tf` into **ALL** corresponding `.tf` files
3. Run `terraform apply` (once for all resources)
4. Validate the state

### Phase 2: Modify All Resources
1. Remove all injected resources from **ALL** `.tf` files
2. Merge `modify.tf` with `add.tf` for each resource (supports partial definitions - only specify changed fields!)
3. Inject the merged resources into **ALL** corresponding `.tf` files
4. Run `terraform apply` (once for all resources)
5. Validate the state

**Important**: `modify.tf` can only be used if a corresponding `add.tf` exists. You cannot modify resources that weren't added first!

### Phase 3: Cleanup All Resources
1. Remove all injected resources from **ALL** `.tf` files
2. Run `terraform apply` (once for cleanup)
3. Restore all original `.tf` files from backups

**Note**: This approach processes all test resources together in batches, which is much faster than processing each resource file individually. The pre-flight check runs once at the start, then all create operations happen, followed by all modifications, and finally all deletions. Each phase only runs `terraform apply` (which does planning automatically), avoiding redundant plan operations.


## Example Output

```
$ cd examples
$ pytest ../test_integration/test_terraform_resources.py -v -s

============================================================
Batch Testing All Resources
============================================================
No test cases found for tenants - skipping
No test cases found for pods - skipping
No test cases found for diagnosticsprofiles - skipping
No test cases found for ipv4prefixlists - skipping
No test cases found for switchpoints - skipping
No test cases found for portacls - skipping
No test cases found for bundles - skipping
No test cases found for packetqueues - skipping
No test cases found for ipv4lists - skipping
No test cases found for ipv6lists - skipping
No test cases found for devicecontrollers - skipping
No test cases found for lags - skipping
No test cases found for diagnosticsportprofiles - skipping
Skipping stages.tf - not in importer.go
No test cases found for policybasedroutingacl - skipping
No test cases found for aspathaccesslists - skipping
No test cases found for spineplanes - skipping
Skipping main.tf - not in importer.go
No test cases found for devicesettings - skipping
No test cases found for acls_ipv6 - skipping
No test cases found for packetbroker - skipping
No test cases found for sfpbreakouts - skipping
No test cases found for routemaps - skipping
No test cases found for ethportsettings - skipping
No test cases found for policybasedrouting - skipping
No test cases found for gatewayprofiles - skipping
No test cases found for services - skipping
No test cases found for extendedcommunitylists - skipping
✓ Found test cases for gateways: ['add', 'modify']
No test cases found for communitylists - skipping
No test cases found for ethportprofiles - skipping
No test cases found for sites - skipping
No test cases found for routemapclauses - skipping
No test cases found for acls_ipv4 - skipping
No test cases found for ipv6prefixlists - skipping
No test cases found for sflowcollectors - skipping

Discovered 1 resources with test cases:
  - gateways (gateways.tf): ['add', 'modify']

Creating backups for all files...
Created backups in /tmp/tf_test_backup_40e44ta5

============================================================
Pre-flight Check: Verifying clean state
============================================================
  Running terraform plan to ensure no pending changes...
  Running terraform plan in /home/terraform-provider-verity/examples...
  ✓ terraform plan successful
  ✓ Clean state confirmed - no pending changes

============================================================
Phase 1: Adding ALL test resources
============================================================
  Injecting resources from test_cases/gateways/add.tf → gateways.tf

  Running terraform apply for all added resources...
  Running terraform apply in /home/terraform-provider-verity/examples...
  ✓ terraform apply successful
  ✓ All resources added successfully (1 files processed)

============================================================
Phase 2: Modifying ALL test resources
============================================================
  Removing previously injected resources...
  Processing modifications for gateways...
  Injecting merged resources into gateways.tf

  Running terraform apply for all modifications...
  Running terraform apply in /home/terraform-provider-verity/examples...
  ✓ terraform apply successful
  ✓ All resources modified successfully (1 files processed)

============================================================
Phase 3: Cleaning up ALL test resources
============================================================
  Removing test resources from all files...
  Running terraform apply for cleanup...
  Running terraform apply in /home/terraform-provider-verity/examples...
  ✓ terraform apply successful
  ✓ All resources cleaned up successfully

  Restoring all files from backup...
  ✓ All files restored from backup

============================================================
✓ All tests completed successfully
  - Total resources tested: 1
  - Resources with add.tf: 1
  - Resources with modify.tf: 1
  - Total terraform applies: 3 (1 per phase)
============================================================

  Final cleanup: Ensuring all files are restored...
PASSED

======================================================== 1 passed in 125.18s (0:02:05)=========================================================
```

