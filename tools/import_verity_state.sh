#!/bin/bash

function log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

log "Verity State Importer Script"
log "============================"

if ! command -v terraform &> /dev/null; then
  log "[ERROR] Terraform command not found. Please install Terraform."
  exit 1
fi

if ! command -v awk &> /dev/null; then
  log "[ERROR] awk command not found. This script requires awk to run properly."
  exit 1
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
  # Create a backup of the original file (which includes the importer)
  cp "$MAIN_TF_FILE" "${MAIN_TF_FILE}.orig"
  
  # Create a clean version without the importer block
  awk 'BEGIN{skip=0} 
       /data "verity_state_importer"/ {skip=1; next} 
       /^[[:space:]]*}/ {if (skip) {skip=0; next}} 
       {if (!skip) print}' "$MAIN_TF_FILE" > "${MAIN_TF_FILE}.clean"
  log "[INFO] Created clean version of $MAIN_TF_FILE without importer"
else
  log "[INFO] verity_state_importer not found, will add it temporarily"
  # Create a backup of the original file (without importer)
  cp "$MAIN_TF_FILE" "${MAIN_TF_FILE}.clean"
  
  # Add the state importer data source
  echo -e "\ndata \"verity_state_importer\" \"import\" {}" >> "$MAIN_TF_FILE"
  log "[INFO] Added verity_state_importer to $MAIN_TF_FILE"
fi

# First terraform apply to generate import blocks
log "[INFO] Running terraform apply to generate resource files and import blocks..."
terraform apply -auto-approve
if [ $? -ne 0 ]; then
  # Restore the original file
  if [ "$IMPORTER_EXISTS" = true ] && [ -f "${MAIN_TF_FILE}.orig" ]; then
    mv "${MAIN_TF_FILE}.orig" "$MAIN_TF_FILE"
  elif [ -f "${MAIN_TF_FILE}.clean" ]; then
    mv "${MAIN_TF_FILE}.clean" "$MAIN_TF_FILE"
  fi
  log "[INFO] Restored original $MAIN_TF_FILE"
  # Clean up backup files
  rm -f "${MAIN_TF_FILE}.orig" "${MAIN_TF_FILE}.clean"
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

# Restore the clean version of the file (without importer)
if [ -f "${MAIN_TF_FILE}.clean" ]; then
  mv "${MAIN_TF_FILE}.clean" "$MAIN_TF_FILE"
  log "[INFO] Removed verity_state_importer from $MAIN_TF_FILE"
fi
rm -f "${MAIN_TF_FILE}.orig"

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
