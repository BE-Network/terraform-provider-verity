# Verity Terraform Provider Documentation

## 1. Provider Configuration

The provider requires basic authentication parameters in your Terraform configuration:

```hcl
terraform {
  required_providers {
    verity = {
      source = "registry.terraform.io/local/verity"
    }
  }
}

provider "verity" {
  uri      = var.uri
  username = var.username
  password = var.password
}
```

Required parameters:
- **uri**: Base URL for the API
- **username** and **password**: For authentication

You can export these variables using environment variables:

```bash
# For Linux/MacOS
export TF_VAR_uri="https://your-verity-instance"
export TF_VAR_username="your-username"
export TF_VAR_password="your-password"

# For Windows PowerShell
$env:TF_VAR_uri="https://your-verity-instance"
$env:TF_VAR_username="your-username"
$env:TF_VAR_password="your-password"
```

## 2. Resource Types

The provider supports the following resource types:

- `verity_tenant`: Manage Verity tenants
- `verity_gateway_profile`: Configure gateway profiles
- `verity_eth_port_profile`: Manage Ethernet port profiles
- `verity_lag`: Configure Link Aggregation Groups
- `verity_service`: Manage services
- `verity_eth_port_settings`: Configure Ethernet port settings
- `verity_bundle`: Manage bundles
- `verity_gateway`: Configure gateways

Each resource type has specific attributes and configurations. See the examples section below for detailed usage.

## 3. State Importer

The provider includes a state importer data source that helps you import existing Verity configurations into your Terraform state:

```hcl
data "verity_state_importer" "import" {
  output_dir = var.config_dir  # defaults to current working directory if not specified
}
```

The state importer workflow:

1. **Configuration Export**: The importer connects to your Verity instance and exports the current configuration
2. **Resource Generation**: It automatically generates Terraform resource files (`.tf`) that map your current Verity configuration to Terraform resources
3. **Import Process**: 
   - For Terraform versions < 1.5: The provider includes bash and PowerShell scripts that use `terraform import` commands to import each resource
   - For Terraform versions ≥ 1.5: Use the provided import scripts to generate `import_blocks.tf` containing import blocks for all resources, allowing you to import everything with a single `terraform apply` command

### Generated Files
The importer generates the following Terraform resource files:
- `bundles.tf` - Bundle resources
- `ethportprofiles.tf` - Ethernet port profile resources
- `ethportsettings.tf` - Ethernet port settings resources
- `gatewayprofiles.tf` - Gateway profile resources
- `gateways.tf` - Gateway resources
- `lags.tf` - LAG resources
- `services.tf` - Service resources
- `tenants.tf` - Tenant resources

Additionally, when using the import scripts:
- For Terraform ≥ 1.5: The scripts will generate `import_blocks.tf` containing import blocks for all resources

## 4. Example Configurations

### Gateway Configuration
```hcl
resource "verity_gateway" "example" {
  name = "example_gateway"
  keepalive_timer = 60
  hold_timer = 180
  gateway_mode = "Dynamic BGP"
  neighbor_as_number = 888
  local_as_number = 222
  ebgp_multihop = 255
  
  enable_bfd = true
  bfd_transmission_interval = 300
  bfd_receive_interval = 300
  bfd_detect_multiplier = 3
  
  static_routes {
    index = 1
    enable = false
    ad_value = 1
  }
}
```

## 5. API Integration

The provider:
- Uses an OpenAPI-generated client for API communication
- Handles authentication via cookie-based sessions
- Supports various API endpoints for different resource types
- Automatically manages state and tracks changes

## 6. Flow of Operations

### Resource Management Flow:
1. Provider Configuration and Authentication
2. Resource Definition in Terraform Files
3. Plan Generation (terraform plan)
4. Change Detection and Validation
5. API Operations Execution (terraform apply)
6. State Update

### Import Flow:
1. Configure State Importer
2. Run Import Operation (generates .tf files)
3. Review Generated Configurations
4. Run Import Scripts (< 1.5) or Use Import Blocks (≥ 1.5)
5. Apply Imported State