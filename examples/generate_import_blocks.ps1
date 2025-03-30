$outputFile = "import_blocks.tf"
$skipFiles = @("provider.tf", "versions.tf", "variables.tf", "terraform.tfvars", "import_blocks.tf")

"# Import blocks for Verity resources" | Out-File -FilePath $outputFile -Encoding utf8
"" | Out-File -FilePath $outputFile -Append -Encoding utf8

$resourceOrder = @{
    "verity_tenant"            = 1
    "verity_service"           = 2
    "verity_eth_port_settings" = 3
    "verity_eth_port_profile"  = 4
    "verity_gateway_profile"   = 5
    "verity_gateway"           = 6
    "verity_lag"               = 7
    "verity_bundle"            = 8
}

$importBlocks = @{}

Write-Output "Processing Terraform files..."

Get-ChildItem -Filter *.tf | ForEach-Object {
    $tfFile = $_.FullName
    if ($skipFiles -contains $_.Name) { return }

    Write-Output "Processing $($_.Name)..."
    $lines = Get-Content $tfFile
    $inBlock = $false
    $braceCount = 0
    $resourceBlock = ""

    foreach ($line in $lines) {
        if (-not $inBlock -and $line -match '^\s*resource\s+"') {
            $inBlock = $true
            $resourceBlock = $line + "`n"
            $open = ([regex]::Matches($line, "{")).Count
            $close = ([regex]::Matches($line, "}")).Count
            $braceCount = $open - $close
            continue
        }
        if ($inBlock) {
            $resourceBlock += $line + "`n"
            $open = ([regex]::Matches($line, "{")).Count
            $close = ([regex]::Matches($line, "}")).Count
            $braceCount += $open - $close
            if ($braceCount -eq 0) {
                if ($resourceBlock -match 'resource\s+"([^"]+)"\s+"([^"]+)"') {
                    $resourceType = $matches[1]
                    $hclName = $matches[2]
                    $nameValue = $hclName
                    if ($resourceBlock -match 'name\s*=\s*"([^"]+)"') {
                        $nameValue = $matches[1]
                    }
                    $importBlock = "import {`n  to = ${resourceType}.${hclName}`n  id = `"$nameValue`"`n}`n`n"
                    if (-not $importBlocks.ContainsKey($resourceType)) {
                        $importBlocks[$resourceType] = ""
                    }
                    $importBlocks[$resourceType] += $importBlock
                }
                $inBlock = $false
                $braceCount = 0
                $resourceBlock = ""
            }
        }
    }
}

$resourceOrder.GetEnumerator() | Sort-Object Value | ForEach-Object {
    $resourceType = $_.Key
    if ($importBlocks.ContainsKey($resourceType)) {
        "# ${resourceType} imports" | Out-File -FilePath $outputFile -Append -Encoding utf8
        $importBlocks[$resourceType] | Out-File -FilePath $outputFile -Append -Encoding utf8
    }
}

Write-Output "Import blocks have been generated in $outputFile"
