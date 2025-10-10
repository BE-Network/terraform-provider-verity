# State Importer Data Source

The `verity_state_importer` data source provides functionality to import existing Verity resources into Terraform state. It automatically fetches resources from the Verity API, creates Terraform configuration files, and generates the appropriate import blocks for each resource.

## Example Usage

```hcl
data "verity_state_importer" "import" {
  # If omitted or empty, the current working directory is used.
  output_dir = "/path/to/dir"
}
```

## Schema

### Optional

- `output_dir` (String) - Directory where the Terraform configuration files will be saved. The directory will be created if it doesn't exist. If not specified or empty, files will be created in the current working directory.

## Files generated

The importer writes multiple `.tf` files into the output directory. The importer always writes `stages.tf` and then may write any of the following files:

- stages.tf
- tenants.tf
- gateways.tf
- gatewayprofiles.tf
- ethportprofiles.tf
- lags.tf
- sflowcollectors.tf
- diagnosticsprofiles.tf
- diagnosticsportprofiles.tf
- policybasedrouting.tf
- services.tf
- ethportsettings.tf
- bundles.tf
- acls_ipv4.tf
- acls_ipv6.tf
- badges.tf
- authenticatedethports.tf
- devicecontrollers.tf
- devicevoicesettings.tf
- packetbroker.tf
- packetqueues.tf
- serviceportprofiles.tf
- voiceportprofiles.tf
- spineplanes.tf
- switchpoints.tf
- aspathaccesslists.tf
- communitylists.tf
- devicesettings.tf
- extendedcommunitylists.tf
- ipv4lists.tf
- ipv4prefixlists.tf
- ipv6lists.tf
- ipv6prefixlists.tf
- routemapclauses.tf
- routemaps.tf
- sfpbreakouts.tf
- sites.tf
- pods.tf
- portacls.tf

Note: Not all files will be created. Tasks are filtered by provider mode and API version compatibility, and the importer skips writing a file if the generated Terraform configuration is empty.

Additionally, the importer writes:
- import_blocks.tf — a generated file containing a sequence of Terraform import blocks for the resources found in the output directory.

## Next Steps

After running the data source:

1. Review the generated configuration files
2. Run `terraform plan` to see what Terraform would change
3. Run `terraform apply` to import the resources into the Terraform state (via import_blocks.tf)

Note: The import operation does not modify the actual resources in Verity; it only creates Terraform configuration to manage those resources.
