# State Importer Data Source

The `verity_state_importer` data source provides functionality to import existing Verity resources into Terraform state. It automatically fetches resources from the Verity API, creates Terraform configuration files, and generates the appropriate import blocks for each resource.

## Example Usage

```hcl
# Create a state importer
data "verity_state_importer" "import" {
  output_dir = var.config_dir
}
```

If the `config_dir` variable is not set, all files will be created in the current working directory.

## Schema

### Required

- `output_dir` (String) - Directory where the Terraform configuration files will be saved. The directory will be created if it doesn't exist. If not specified or empty, files will be created in the current working directory.

### Read-Only

- `id` (String) - Identifier for this import operation (timestamp of when the import was executed)
- `imported_files` (List of String) - List of paths to files that were created during the import process

## Generated Files

The data source will generate several Terraform configuration files:

- `tenants.tf` - Verity tenant resources
- `services.tf` - Verity service resources
- `ethportsettings.tf` - Verity Ethernet port settings resources
- `ethportprofiles.tf` - Verity Ethernet port profile resources
- `gatewayprofiles.tf` - Verity gateway profile resources
- `gateways.tf` - Verity gateway resources
- `lags.tf` - Verity LAG resources
- `bundles.tf` - Verity bundle resources
- `import_blocks.tf` - Terraform import blocks for all resources

## Import Order

Resources are imported in a specific order to respect dependencies:

1. `verity_tenant`
2. `verity_service`
3. `verity_eth_port_settings`
4. `verity_eth_port_profile`
5. `verity_gateway_profile`
6. `verity_gateway`
7. `verity_lag`
8. `verity_bundle`

## Next Steps

After running the data source:

1. Review the generated configuration files
2. Run `terraform plan` to see what Terraform would change
3. Run `terraform apply` to import the resources into the Terraform state

Note: The import operation does not modify the actual resources in Verity; it only creates Terraform configuration to manage those resources.
