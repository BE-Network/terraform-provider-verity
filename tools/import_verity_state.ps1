#Requires -Version 5.0

function Log {
    param ([string]$message, [string]$color = "White")
    Write-Host "[$([DateTime]::Now.ToString('yyyy-MM-dd HH:mm:ss'))] $message" -ForegroundColor $color
}

Log "Verity State Importer Script" -color Cyan
Log "============================" -color Cyan

if ($Local) {
    Log "[INFO] Running in local provider mode - will skip terraform init" -color Yellow
}

if (-not (Get-Command terraform -ErrorAction SilentlyContinue)) {
    Log "[ERROR] Terraform command not found. Please install Terraform." -color Red
    exit 1
}

# Only run terraform init if not in local mode
if (-not $Local) {
    Log "[INFO] Running terraform init..." -color Cyan
    terraform init
    if ($LASTEXITCODE -ne 0) {
        Log "[ERROR] Terraform init failed. Exiting." -color Red
        exit 1
    }
}

# Find the main terraform file with the Verity provider
Log "[INFO] Finding main Terraform file with Verity provider..." -color Cyan
$mainTfFile = $null
$tfFiles = Get-ChildItem -Filter "*.tf" -ErrorAction SilentlyContinue
if ($tfFiles.Count -eq 0) {
    Log "[ERROR] No .tf files found in the current directory." -color Red
    exit 1
}

foreach ($file in $tfFiles) {
    $content = Get-Content $file.FullName -Raw
    if ($content -match 'provider\s+"verity"') {
        $mainTfFile = $file.FullName
        Log "[INFO] Found Verity provider in $($file.Name)" -color Cyan
        break
    }
}

if (-not $mainTfFile) {
    Log "[ERROR] No Terraform file with Verity provider found. Cannot continue." -color Red
    exit 1
}

# Check if state importer exists in the file
$importerExists = $false
$fileContent = Get-Content $mainTfFile -Raw
if ($fileContent -match 'verity_state_importer') {
    $importerExists = $true
    Log "[INFO] Found verity_state_importer in $($mainTfFile)" -color Cyan
} else {
    Log "[INFO] verity_state_importer not found, will add it temporarily" -color Cyan
    # Create a backup of the original file
    Copy-Item -Path $mainTfFile -Destination "$mainTfFile.bak" -Force
    
    # Add the state importer data source
    $importerCode = "`ndata `"verity_state_importer`" `"import`" {`n  output_dir = var.config_dir`n}`n"
    Add-Content -Path $mainTfFile -Value $importerCode
    Log "[INFO] Added verity_state_importer to $($mainTfFile)" -color Cyan
}

# First terraform apply to generate import blocks
Log "[INFO] Running terraform apply to generate resource files and import blocks..." -color Cyan
terraform apply -auto-approve
if ($LASTEXITCODE -ne 0) {
    # If we added the importer and apply failed, restore the original file
    if (-not $importerExists -and (Test-Path "$mainTfFile.bak")) {
        Copy-Item -Path "$mainTfFile.bak" -Destination $mainTfFile -Force
        Log "[INFO] Restored original $($mainTfFile)" -color Yellow
    }
    Log "[ERROR] First terraform apply failed. Exiting." -color Red
    exit 1
}

# Check if import_blocks.tf was generated
if (-not (Test-Path "import_blocks.tf")) {
    Log "[ERROR] import_blocks.tf was not generated. Check your configuration." -color Red
    exit 1
} else {
    Log "[INFO] Successfully generated import_blocks.tf" -color Green
}

# If we added the importer, remove it now as it's not needed for the import process
if (-not $importerExists) {
    if (Test-Path "$mainTfFile.bak") {
        Copy-Item -Path "$mainTfFile.bak" -Destination $mainTfFile -Force
        Remove-Item -Path "$mainTfFile.bak" -Force
        Log "[INFO] Removed verity_state_importer from $($mainTfFile)" -color Cyan
    }
}

# Second terraform apply to import resources
Log "[INFO] Running second terraform apply to import resources into state..." -color Cyan
terraform apply -auto-approve
if ($LASTEXITCODE -ne 0) {
    Log "[ERROR] Resource import failed." -color Red
    exit 1
}

# Remove import_blocks.tf if import was successful
Log "[INFO] Import successful. Removing import_blocks.tf..." -color Green
Remove-Item -Path "import_blocks.tf" -Force
Log "[INFO] Import process completed successfully!" -color Green
