# State Importer Data Source

The `verity_state_importer` data source provides functionality to import existing Verity resources into Terraform state. It automatically fetches resources from the Verity API, creates Terraform configuration files, and generates the appropriate import blocks for each resource.

## Example Usage

```hcl
data "verity_state_importer" "import" {
  output_dir = "/path/to/dir" # defaults to current working directory if not specified
}
```

## Schema

### Optional

- `output_dir` (String) - Directory where the Terraform configuration files will be saved. The directory will be created if it doesn't exist. If not specified or empty, files will be created in the current working directory.

## Files generated

The importer writes multiple `.tf` files into the output directory. The importer always writes `stages.tf` and then may write any of the following files:

- `deviceaaaprofiles.tf`
- `acls_ipv4.tf`
- `acls_ipv6.tf`
- `aspathaccesslists.tf`
- `authenticatedethports.tf`
- `badges.tf`
- `bundles.tf`
- `communitylists.tf`
- `devicesettings.tf`
- `devicevoicesettings.tf`
- `diagnosticsportprofiles.tf`
- `diagnosticsprofiles.tf`
- `ethportprofiles.tf`
- `ethportsettings.tf`
- `extendedcommunitylists.tf`
- `gatewayprofiles.tf`
- `gateways.tf`
- `groupingrules.tf`
- `ipv4lists.tf`
- `ipv4prefixlists.tf`
- `ipv6lists.tf`
- `ipv6prefixlists.tf`
- `lags.tf`
- `macfilters.tf`
- `packetbroker.tf`
- `packetqueues.tf`
- `pairs.tf`
- `policybasedrouting.tf`
- `policybasedroutingacl.tf`
- `pods.tf`
- `portacls.tf`
- `routemapclauses.tf`
- `routemaps.tf`
- `services.tf`
- `serviceportprofiles.tf`
- `sflowcollectors.tf`
- `sfpbreakouts.tf`
- `sites.tf`
- `spineplanes.tf`
- `sspgroups.tf`
- `switchpoints.tf`
- `sus.tf`
- `tacacsprofiles.tf`
- `tenants.tf`
- `thresholdgroups.tf`
- `thresholds.tf`
- `voiceportprofiles.tf`
- `import_blocks.tf` - Import blocks for all resources in the correct dependency order
- `stages.tf` - Resource dependency ordering

Note: Not all files will be created. Tasks are filtered by provider mode and API version compatibility, and the importer skips writing a file if the generated Terraform configuration is empty.

Additionally, the importer writes:
- import_blocks.tf — a generated file containing a sequence of Terraform import blocks for the resources found in the output directory.

## Next Steps

After running the data source:

1. Review the generated configuration files
2. Run `terraform plan` to see what Terraform would change
3. Run `terraform apply` to import the resources into the Terraform state (via import_blocks.tf)

Note: The import operation does not modify the actual resources in Verity; it only creates Terraform configuration to manage those resources.
