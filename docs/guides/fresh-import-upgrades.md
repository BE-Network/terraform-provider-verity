# Upgrading with fresh configuration and state

Breaking API releases can require a fresh import instead of upgrading existing
Terraform state. Check the target provider's release upgrade guide for removed
endpoints, changed fields, and objects it can no longer manage. This procedure
adopts the current server configuration using the target provider.

1. Archive the existing Terraform configuration, provider version, and state.
   If using a remote backend, preserve its state snapshot too. Unapplied changes
   in the old configuration do not transfer automatically.
2. Prepare a fresh Terraform directory containing only the provider and backend
   configuration needed for the target release. Pin the target provider version
   and select the correct server and mode. See the
   [provider setup](../../README.md#production-setup). Use empty Terraform state:
   a new local directory is insufficient if its backend still loads the old
   state. Retire the old workspace's automation before using the replacement.
3. Run `terraform init` in that directory, then run the import script from the
   same directory:

   ```bash
   bash /path/to/terraform-provider-verity/tools/import_verity_state.sh
   ```

   On Windows:

   ```powershell
   & "C:\path\to\terraform-provider-verity\tools\import_verity_state.ps1"
   ```

   The scripts require an existing `.tf` file with the Verity provider
   configuration. They run two Terraform applies: first to generate current
   resource configuration and import blocks, then to import those objects into
   the empty state. They do not reset an old backend or migrate old state.
4. Review the generated configuration, imported state, and any
   `unsupported_arguments.txt` output. Run `terraform plan` and investigate
   unexpected changes before modifying resources. Complete any manual
   configuration required by the release guide.

Objects behind removed endpoints may still exist on the server but are no longer
managed by this provider. Their omission from the new import does not delete them.
Retain the archived configuration and state for reference; use the new workspace
for subsequent changes.

Maintainers must describe release-specific removals, replacements, field changes,
and manual configuration in an upgrade guide referenced by
`specs/state_transitions.json`. Refreshing a schema snapshot alone does not approve
a breaking release. See the
[release compatibility process](../../README.md#breaking-api-changes-and-existing-terraform-state).
