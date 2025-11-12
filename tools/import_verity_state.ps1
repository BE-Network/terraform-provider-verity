#Requires -Version 5.0

function Log {
    param ([string]$message, [string]$color = "White")
    Write-Host "[$([DateTime]::Now.ToString('yyyy-MM-dd HH:mm:ss'))] $message" -ForegroundColor $color
}

Log "Verity State Importer Script" -color Cyan
Log "============================" -color Cyan

if (-not (Get-Command terraform -ErrorAction SilentlyContinue)) {
    Log "[ERROR] Terraform command not found. Please install Terraform." -color Red
    exit 1
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
    Copy-Item -Path $mainTfFile -Destination "$mainTfFile.orig" -Force
    
    # Create a clean version without the importer block
    $content = Get-Content $mainTfFile
    $cleanContent = New-Object System.Collections.ArrayList
    $skipLines = $false
    
    foreach ($line in $content) {
        # Start skipping at the data "verity_state_importer" line
        if ($line -match 'data\s+"verity_state_importer"') {
            $skipLines = $true
            continue
        }
        
        # If we find a closing brace that completes the block, stop skipping after this line
        if ($skipLines -and $line -match '^\s*}\s*$') {
            $skipLines = $false
            continue
        }
        
        if (-not $skipLines) {
            [void]$cleanContent.Add($line)
        }
    }
    
    $cleanContent | Set-Content -Path "$mainTfFile.clean" -Force
} else {
    Log "[INFO] verity_state_importer not found, will add it temporarily" -color Cyan
    # Create a backup of the original file (without importer)
    Copy-Item -Path $mainTfFile -Destination "$mainTfFile.clean" -Force
    
    # Add the state importer data source
    $importerCode = "`ndata `"verity_state_importer`" `"import`" {}`n"
    Add-Content -Path $mainTfFile -Value $importerCode
    Log "[INFO] Added verity_state_importer to $($mainTfFile)" -color Cyan
}

# First terraform apply to generate import blocks
Log "[INFO] Running terraform apply to generate resource files and import blocks..." -color Cyan
terraform apply -auto-approve
if ($LASTEXITCODE -ne 0) {
    # Restore the original file
    if ($importerExists -and (Test-Path "$mainTfFile.orig")) {
        Copy-Item -Path "$mainTfFile.orig" -Destination $mainTfFile -Force
        Log "[INFO] Restored original file with importer" -color Yellow
    } elseif (Test-Path "$mainTfFile.clean") {
        Copy-Item -Path "$mainTfFile.clean" -Destination $mainTfFile -Force
        Log "[INFO] Restored original file without importer" -color Yellow
    }

    Remove-Item -Path "$mainTfFile.orig" -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "$mainTfFile.clean" -Force -ErrorAction SilentlyContinue
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

# Always restore the clean version of the file (without importer)
if (Test-Path "$mainTfFile.clean") {
    Copy-Item -Path "$mainTfFile.clean" -Destination $mainTfFile -Force
    Remove-Item -Path "$mainTfFile.clean" -Force
    Log "[INFO] Removed verity_state_importer from $($mainTfFile)" -color Cyan
}

Remove-Item -Path "$mainTfFile.orig" -Force -ErrorAction SilentlyContinue

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
