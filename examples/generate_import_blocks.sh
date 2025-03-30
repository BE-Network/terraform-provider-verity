#!/bin/bash

skip_files=("provider.tf" "versions.tf" "variables.tf" "terraform.tfvars" "import_blocks.tf")
output_file="import_blocks.tf"

{
  echo "# Import blocks for Verity resources"
  echo ""
} > "$output_file"

resource_order=( verity_tenant verity_service verity_eth_port_settings verity_eth_port_profile verity_gateway_profile verity_gateway verity_lag verity_bundle )

import_blocks_verity_tenant=""
import_blocks_verity_service=""
import_blocks_verity_eth_port_settings=""
import_blocks_verity_eth_port_profile=""
import_blocks_verity_gateway_profile=""
import_blocks_verity_gateway=""
import_blocks_verity_lag=""
import_blocks_verity_bundle=""

echo "Processing Terraform files..."
for tf_file in *.tf; do
  skip=no
  for skip_file in "${skip_files[@]}"; do
    if [[ "$tf_file" == "$skip_file" ]]; then
      skip=yes
      break
    fi
  done
  [[ "$skip" == "yes" ]] && continue

  echo "Processing $tf_file..."
  in_block=0
  resource_block=""
  
  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ $in_block -eq 0 && $line =~ ^[[:space:]]*resource[[:space:]]+\" ]]; then
      in_block=1
      resource_block="$line"$'\n'
      open_braces=$(grep -o "{" <<< "$line" | wc -l)
      close_braces=$(grep -o "}" <<< "$line" | wc -l)
      brace_count=$((open_braces - close_braces))
      continue
    fi

    if [[ $in_block -eq 1 ]]; then
      resource_block+="$line"$'\n'
      open_braces=$(grep -o "{" <<< "$line" | wc -l)
      close_braces=$(grep -o "}" <<< "$line" | wc -l)
      brace_count=$((brace_count + open_braces - close_braces))
      if [[ $brace_count -eq 0 ]]; then
        if [[ $resource_block =~ resource[[:space:]]+\"([^\"]+)\"[[:space:]]+\"([^\"]+)\" ]]; then
          resource_type="${BASH_REMATCH[1]}"
          hcl_name="${BASH_REMATCH[2]}"
          if [[ $resource_block =~ name[[:space:]]*=[[:space:]]*\"([^\"]+)\" ]]; then
            name_value="${BASH_REMATCH[1]}"
          else
            name_value="$hcl_name"
          fi
  
          import_block="import {
  to = ${resource_type}.${hcl_name}
  id = \"${name_value}\"
}

"
          case "$resource_type" in
            verity_tenant)
              import_blocks_verity_tenant="${import_blocks_verity_tenant}${import_block}"
              ;;
            verity_service)
              import_blocks_verity_service="${import_blocks_verity_service}${import_block}"
              ;;
            verity_eth_port_settings)
              import_blocks_verity_eth_port_settings="${import_blocks_verity_eth_port_settings}${import_block}"
              ;;
            verity_eth_port_profile)
              import_blocks_verity_eth_port_profile="${import_blocks_verity_eth_port_profile}${import_block}"
              ;;
            verity_gateway_profile)
              import_blocks_verity_gateway_profile="${import_blocks_verity_gateway_profile}${import_block}"
              ;;
            verity_gateway)
              import_blocks_verity_gateway="${import_blocks_verity_gateway}${import_block}"
              ;;
            verity_lag)
              import_blocks_verity_lag="${import_blocks_verity_lag}${import_block}"
              ;;
            verity_bundle)
              import_blocks_verity_bundle="${import_blocks_verity_bundle}${import_block}"
              ;;
            *)
              echo "Unknown resource type: $resource_type"
              ;;
          esac
        fi
        in_block=0
        resource_block=""
      fi
    fi
  done < "$tf_file"
done

for res in "${resource_order[@]}"; do
  eval "block=\$import_blocks_${res}"
  if [[ -n "$block" ]]; then
    {
      echo "# ${res} imports"
      echo "$block"
    } >> "$output_file"
  fi
done

echo "Import blocks have been generated in $output_file"
