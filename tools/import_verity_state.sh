#!/bin/bash

function log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

log "Verity State Importer Script"
log "============================"

LOCAL_MODE=false
for arg in "$@"; do
  if [ "$arg" == "--local" ]; then
    LOCAL_MODE=true
    log "[INFO] Running in local provider mode - will skip terraform init"
  fi
done

if ! command -v terraform &> /dev/null; then
  log "[ERROR] Terraform command not found. Please install Terraform."
  exit 1
fi

# Only run terraform init if not in local mode
if [ "$LOCAL_MODE" = false ]; then
  log "[INFO] Running terraform init..."
  terraform init
  if [ $? -ne 0 ]; then
    log "[ERROR] Terraform init failed. Exiting."
    exit 1
  fi
fi

# Find the main terraform file with the Verity provider
log "[INFO] Finding main Terraform file with Verity provider..."
MAIN_TF_FILE=""
for file in *.tf; do
  if [ -f "$file" ] && grep -q "provider \"verity\"" "$file"; then
    MAIN_TF_FILE="$file"
    log "[INFO] Found Verity provider in $MAIN_TF_FILE"
    break
  fi
done

if [ -z "$MAIN_TF_FILE" ]; then
  log "[ERROR] No Terraform file with Verity provider found. Cannot continue."
  exit 1
fi

# Check if state importer exists in the file
IMPORTER_EXISTS=false
if grep -q "verity_state_importer" "$MAIN_TF_FILE"; then
  IMPORTER_EXISTS=true
  log "[INFO] Found verity_state_importer in $MAIN_TF_FILE"
else
  log "[INFO] verity_state_importer not found, will add it temporarily"
  # Create a backup of the original file
  cp "$MAIN_TF_FILE" "${MAIN_TF_FILE}.bak"
  
  # Add the state importer data source
  echo -e "\ndata \"verity_state_importer\" \"import\" {\n  output_dir = var.config_dir\n}" >> "$MAIN_TF_FILE"
  log "[INFO] Added verity_state_importer to $MAIN_TF_FILE"
fi

# First terraform apply to generate import blocks
log "[INFO] Running terraform apply to generate resource files and import blocks..."
terraform apply -auto-approve
if [ $? -ne 0 ]; then
  # If we added the importer and apply failed, restore the original file
  if [ "$IMPORTER_EXISTS" = false ] && [ -f "${MAIN_TF_FILE}.bak" ]; then
    mv "${MAIN_TF_FILE}.bak" "$MAIN_TF_FILE"
    log "[INFO] Restored original $MAIN_TF_FILE"
  fi
  log "[ERROR] First terraform apply failed. Exiting."
  exit 1
fi

# Check if import_blocks.tf was generated
if [ ! -f "import_blocks.tf" ]; then
  log "[ERROR] import_blocks.tf was not generated. Check your configuration."
  exit 1
else
  log "[INFO] Successfully generated import_blocks.tf"
fi

# If we added the importer, remove it now as it's not needed for the import process
if [ "$IMPORTER_EXISTS" = false ]; then
  if [ -f "${MAIN_TF_FILE}.bak" ]; then
    mv "${MAIN_TF_FILE}.bak" "$MAIN_TF_FILE"
    log "[INFO] Removed verity_state_importer from $MAIN_TF_FILE"
  fi
fi

# Second terraform apply to import resources
log "[INFO] Running second terraform apply to import resources into state..."
terraform apply -auto-approve
if [ $? -ne 0 ]; then
  log "[ERROR] Resource import failed."
  exit 1
fi

# Remove import_blocks.tf if import was successful
log "[INFO] Import successful. Removing import_blocks.tf..."
rm -f import_blocks.tf
log "[INFO] Import process completed successfully!"
