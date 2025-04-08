# Local Verity Terraform Provider Setup

## Building the Provider

To compile the provider binary, run the following command in the root directory of the project:

```bash
go build -o terraform-provider-verity
```

If you are using windows...

```bash
go build -o terraform-provider-verity.exe
```


## Provider Configuration

### Custom Provider Binary

To use a local development version of the provider, you need to configure Terraform to use your custom provider binary instead of downloading it from the registry. This is done using a `.tfrc` configuration file. Here's an example (`dev.tfrc`):

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/local/verity" = "/home/<user>/terraform-provider-verity"
  }
  direct {}
}
```
(Note: if using windows pay attention to the windows format for director naming)


To use this configuration, you can either:
- Copy it to `~/.terraformrc`
- Set the `TF_CLI_CONFIG_FILE` environment variable to point to your custom `.tfrc` file:
  ```bash
  export TF_CLI_CONFIG_FILE=/home/<user>/terraform-provider-verity/examples/dev.tfrc
  ```
  For Windows:
  ```powershell
  $env:TF_CLI_CONFIG_FILE="C:\path\to\terraform-provider-verity\examples\dev.tfrc"
  ```

### Required Environment Variables

The provider requires several environment variables to be set for authentication. You can set these variables manually or use a script like `set_env.sh`:

```bash
export TF_VAR_uri="<your-verity-uri>"
export TF_VAR_username="<your-username>"
export TF_VAR_password="<your-password>"
```

For Windows, use the following commands to set environment variables:

```powershell
$env:TF_VAR_uri="<your-verity-uri>"
$env:TF_VAR_username="<your-username>"
$env:TF_VAR_password="<your-password>"
```

These variables correspond to the provider configuration in your Terraform files:

```hcl
provider "verity" {
  uri      = var.uri
  username = var.username
  password = var.password
}
```

Make sure to set these environment variables before running any Terraform commands. (Of course, windows is different...)


## Recommended Environment Variables

By default, Terraform is only able to handle 10 resources at a time. It is recommended to export the following environment variable to allow the provider to make bulk requests to the API instead of splitting the API calls into multiple ones:

### Unix-based Systems
```bash
export TF_CLI_ARGS_apply="-parallelism=1000"
```

### Windows
```powershell
$env:TF_CLI_ARGS_apply="-parallelism=1000"
```

## Regenerating the OpenAPI Go SDK

If there are changes to the Verity API (reflected in a new swagger.json file), you'll need to regenerate the OpenAPI Go SDK. Follow these steps:

1. Save the new swagger.json file

2. Run the tools/transform_swagger.py script to prepare the swagger file for code generation:
   ```bash
   python3 tools/transform_swagger.py
   ```

3. Install the OpenAPI Generator CLI (if not already installed):
   ```bash
   npm install @openapitools/openapi-generator-cli -g
   ```

4. Generate the Go SDK using openapi-generator-cli:
   ```bash
   openapi-generator-cli generate -i swagger_transformed.json -g go -o ./openapi
   ```

5. Replace the existing openapi folder with the newly generated one

### Updating Provider Resource Files

After regenerating the SDK, you need to update the provider resource files:

- For fields deleted from the API: Remove them from the corresponding provider resource files
- For new fields added to the API: Add them to the appropriate provider resource files
