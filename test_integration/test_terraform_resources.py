#!/usr/bin/env python3
import sys
import shutil
import subprocess
import json
import tempfile
import re
from pathlib import Path
from typing import Dict, List, Tuple, Optional, Set, Any

# Check for required dependencies before proceeding
try:
    import pytest
except ImportError:
    print("❌ ERROR: pytest is not installed")
    print("   Install it with: pip install -r requirements.txt")
    print("   Or manually: pip install pytest")
    sys.exit(1)

try:
    import hcl2
except ImportError:
    print("❌ ERROR: python-hcl2 is not installed")
    print("   This is required for partial resource definitions in modify.tf")
    print("   Install it with: pip install -r requirements.txt")
    print("   Or manually: pip install python-hcl2")
    sys.exit(1)


class TerraformTestRunner:
    """Handles terraform operations and test execution."""

    def __init__(self, tf_dir: Optional[str] = None, test_cases_dir: Optional[str] = None):
        # Use current directory if not specified
        self.tf_dir = Path(tf_dir) if tf_dir else Path.cwd()

        # Find test_cases directory relative to script location
        script_dir = Path(__file__).parent
        if test_cases_dir:
            self.test_cases_dir = Path(test_cases_dir)
        else:
            self.test_cases_dir = script_dir / "test_cases"

        self.backup_dir = None
        self.valid_resources = self._load_valid_resources()

    def _load_valid_resources(self) -> Set[str]:
        """Load valid resource names from importer.go."""
        valid_resources = set()

        # Find importer.go relative to script location
        script_dir = Path(__file__).parent
        importer_file = script_dir.parent / "internal" / "importer" / "importer.go"

        if not importer_file.exists():
            print(f"⚠ Warning: importer.go not found at {importer_file}")
            print(f"  All .tf files will be tested")
            return valid_resources

        try:
            with open(importer_file, 'r') as f:
                content = f.read()

            pattern = r'\{name:\s*"([^"]+)",'
            matches = re.findall(pattern, content)
            valid_resources = set(matches)

            print(f"Loaded {len(valid_resources)} valid resources from importer.go")
            return valid_resources

        except Exception as e:
            print(f"Warning: Failed to load importer.go: {e}")
            print(f"All .tf files will be tested")
            return valid_resources

    def run_terraform_command(self, command: List[str], cwd: str) -> Tuple[int, str, str]:
        """Run a terraform command and return exit code, stdout, stderr."""
        try:
            result = subprocess.run(
                command,
                cwd=cwd,
                capture_output=True,
                text=True,
                timeout=300
            )
            return result.returncode, result.stdout, result.stderr
        except subprocess.TimeoutExpired:
            return -1, "", "Command timed out after 300 seconds"
        except Exception as e:
            return -1, "", str(e)

    def terraform_plan(self, cwd: str) -> Tuple[bool, str]:
        """Run terraform plan and return success status and output."""
        print(f"  Running terraform plan in {cwd}...")
        returncode, stdout, stderr = self.run_terraform_command(
            ["terraform", "plan", "-out=tfplan"],
            cwd
        )
        if returncode != 0:
            print(f"  ❌ terraform plan failed:")
            print(f"    stderr: {stderr}")
            return False, stderr
        print(f"  ✓ terraform plan successful")
        return True, stdout

    def terraform_apply(self, cwd: str) -> bool:
        """Apply terraform changes."""
        print(f"  Running terraform apply in {cwd}...")
        returncode, stdout, stderr = self.run_terraform_command(
            ["terraform", "apply", "-auto-approve", "tfplan"],
            cwd
        )
        if returncode != 0:
            print(f"  ❌ terraform apply failed:")
            print(f"    stdout: {stdout}")
            print(f"    stderr: {stderr}")
            return False
        print(f"  ✓ terraform apply successful")
        
        tfplan_path = Path(cwd) / "tfplan"
        if tfplan_path.exists():
            try:
                tfplan_path.unlink()
                print(f"  ✓ Cleaned up tfplan file")
            except Exception as e:
                print(f"  ⚠ Warning: Failed to remove tfplan file: {e}")
        
        return True

    def terraform_show(self, cwd: str) -> Optional[Dict]:
        """Get current terraform state as JSON."""
        returncode, stdout, stderr = self.run_terraform_command(
            ["terraform", "show", "-json"],
            cwd
        )
        if returncode != 0:
            print(f"  ⚠ terraform show failed: {stderr}")
            return None
        try:
            return json.loads(stdout)
        except json.JSONDecodeError as e:
            print(f"  ⚠ Failed to parse terraform state: {e}")
            return None

    def create_backup(self, file_path: Path) -> Path:
        """Create a backup of the file."""
        if self.backup_dir is None:
            self.backup_dir = Path(tempfile.mkdtemp(prefix="tf_test_backup_"))

        backup_path = self.backup_dir / file_path.name
        shutil.copy2(file_path, backup_path)
        return backup_path

    def restore_backup(self, original_path: Path, backup_path: Path):
        """Restore a file from backup."""
        shutil.copy2(backup_path, original_path)

    def inject_resources(self, target_file: Path, test_case_file: Path) -> bool:
        """Inject test resources into the target file."""
        try:
            with open(test_case_file, 'r') as f:
                test_content = f.read()

            with open(target_file, 'a') as f:
                f.write("\n\n# === TEST RESOURCES START ===\n")
                f.write(test_content)
                f.write("\n# === TEST RESOURCES END ===\n")

            return True
        except Exception as e:
            print(f"  ❌ Failed to inject resources: {e}")
            return False

    def remove_injected_resources(self, target_file: Path):
        """Remove injected test resources from the target file."""
        try:
            with open(target_file, 'r') as f:
                content = f.read()

            # Remove everything between markers
            start_marker = "# === TEST RESOURCES START ==="
            end_marker = "# === TEST RESOURCES END ==="

            if start_marker in content:
                start_idx = content.find(start_marker)
                end_idx = content.find(end_marker)
                if end_idx != -1:
                    # Remove including the end marker and trailing newline
                    new_content = content[:start_idx] + content[end_idx + len(end_marker) + 1:]
                    # Clean up any trailing whitespace
                    new_content = new_content.rstrip() + "\n" if new_content.strip() else ""

                    with open(target_file, 'w') as f:
                        f.write(new_content)

            return True
        except Exception as e:
            print(f"  ❌ Failed to remove injected resources: {e}")
            return False

    def _merge_hcl_attributes(self, base: Dict[str, Any], override: Dict[str, Any]) -> Dict[str, Any]:
        """Merge HCL attributes, with override taking precedence."""
        merged = base.copy()
        
        for key, value in override.items():
            if key in merged and isinstance(merged[key], dict) and isinstance(value, dict):
                # Recursively merge nested dictionaries
                merged[key] = self._merge_hcl_attributes(merged[key], value)
            else:
                # Override the value
                merged[key] = value
        
        return merged

    def _format_hcl_value(self, value: Any, indent: int = 0) -> str:
        """Format a Python value as HCL syntax."""
        indent_str = "  " * indent
        
        if value is None:
            return "null"
        elif isinstance(value, bool):
            return "true" if value else "false"
        elif isinstance(value, str):
            # Check if it's a terraform interpolation wrapped by hcl2 parser (${...})
            # These should be unwrapped for HCL output (e.g., depends_on)
            if value.startswith('${') and value.endswith('}'):
                # Strip the ${ } wrapper added by hcl2 parser
                return value[2:-1]
            # Check if it's a reference (e.g., verity_gateway.gateway_1.id)
            if re.match(r'^[a-zA-Z_][a-zA-Z0-9_]*\.[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)*$', value):
                return value
            # Regular string
            return f'"{value}"'
        elif isinstance(value, (int, float)):
            return str(value)
        elif isinstance(value, list):
            if not value:
                return "[]"
            # Format list items
            items = [self._format_hcl_value(item, indent) for item in value]
            if all(isinstance(item, (str, int, float, bool)) or item is None for item in value):
                # Simple list on one line
                return "[" + ", ".join(items) + "]"
            else:
                # Complex list with line breaks
                formatted_items = [f"\n{indent_str}  {item}" for item in items]
                return "[" + ",".join(formatted_items) + f"\n{indent_str}]"
        elif isinstance(value, dict):
            if not value:
                return "{}"
            # Format dictionary as HCL block
            lines = ["{"]
            for k, v in value.items():
                formatted_value = self._format_hcl_value(v, indent + 1)
                lines.append(f"{indent_str}  {k} = {formatted_value}")
            lines.append(f"{indent_str}}}")
            return "\n".join(lines)
        else:
            return str(value)

    def _format_hcl_resource(self, resource_type: str, resource_name: str, attributes: Dict[str, Any]) -> str:
        """Format a resource as HCL syntax."""
        lines = [f'resource "{resource_type}" "{resource_name}" {{']
        
        for key, value in attributes.items():
            # Check if value is a list of dicts (HCL blocks)
            if isinstance(value, list) and value and all(isinstance(item, dict) for item in value):
                # Format as HCL blocks
                for block in value:
                    lines.append(f"  {key} {{")
                    for k, v in block.items():
                        formatted_value = self._format_hcl_value(v, indent=2)
                        lines.append(f"    {k} = {formatted_value}")
                    lines.append("  }")
            else:
                formatted_value = self._format_hcl_value(value, indent=1)
                lines.append(f"  {key} = {formatted_value}")
        
        lines.append("}")
        return "\n".join(lines)

    def merge_modify_with_add(self, add_file: Path, modify_file: Path) -> Optional[str]:
        """
        Merge modify.tf with add.tf, allowing partial resource definitions in modify.tf.
        Returns the merged HCL content as a string.
        Returns None if files are empty or parsing fails.
        """
        try:
            # Parse both files
            with open(add_file, 'r') as f:
                add_content = hcl2.load(f)
            
            with open(modify_file, 'r') as f:
                modify_content = hcl2.load(f)
            
            # Extract resources from add.tf
            add_resources = {}
            if 'resource' in add_content:
                for resource_list in add_content['resource']:
                    for resource_type, resources in resource_list.items():
                        for resource_name, attributes in resources.items():
                            key = (resource_type, resource_name)
                            add_resources[key] = attributes
            
            # Check if add.tf is empty (no resources)
            if not add_resources:
                print(f"  ⚠️  WARNING: add.tf has no resources defined (only comments)")
                return None
            
            # Extract resources from modify.tf and merge
            merged_resources = add_resources.copy()
            if 'resource' in modify_content:
                for resource_list in modify_content['resource']:
                    for resource_type, resources in resource_list.items():
                        for resource_name, attributes in resources.items():
                            key = (resource_type, resource_name)
                            if key in merged_resources:
                                # Merge with existing resource
                                merged_resources[key] = self._merge_hcl_attributes(
                                    merged_resources[key],
                                    attributes
                                )
                            else:
                                # New resource in modify.tf
                                merged_resources[key] = attributes
            
            # Format back to HCL
            output_lines = []
            for (resource_type, resource_name), attributes in merged_resources.items():
                hcl_resource = self._format_hcl_resource(resource_type, resource_name, attributes)
                output_lines.append(hcl_resource)
                output_lines.append("")  # Empty line between resources
            
            return "\n".join(output_lines)
            
        except Exception as e:
            print(f"  ❌ Failed to merge HCL files: {e}")
            print(f"  Please check your add.tf and modify.tf syntax")
            return None

    def find_test_cases(self, resource_name: str) -> Dict[str, Path]:
        """Find test case files for a given resource and validate they contain resources."""
        test_cases = {}
        resource_test_dir = self.test_cases_dir / resource_name

        if not resource_test_dir.exists():
            return test_cases

        # Look for add.tf and modify.tf
        add_file = resource_test_dir / "add.tf"
        if add_file.exists():
            # Validate add.tf has actual resources
            if self._file_has_resources(add_file):
                test_cases["add"] = add_file

        modify_file = resource_test_dir / "modify.tf"
        if modify_file.exists():
            # Validate modify.tf has actual resources
            if self._file_has_resources(modify_file):
                test_cases["modify"] = modify_file

        return test_cases

    def _file_has_resources(self, file_path: Path) -> bool:
        """Check if a .tf file contains actual resource definitions (not just comments)."""
        try:
            with open(file_path, 'r') as f:
                content = hcl2.load(f)
            
            # Check if file has any resources
            if 'resource' in content and content['resource']:
                return True
            
            return False
        except Exception:
            # If parsing fails, assume file is invalid/empty
            return False

    def cleanup(self):
        """Clean up backup directory."""
        if self.backup_dir and self.backup_dir.exists():
            shutil.rmtree(self.backup_dir)


class TestTerraformResources:
    """Pytest test class for terraform resources."""

    @pytest.fixture(scope="class")
    def runner(self):
        """Create a test runner instance."""
        runner = TerraformTestRunner()
        yield runner
        runner.cleanup()

    def discover_tf_files_with_test_cases(self, runner: TerraformTestRunner) -> List[Tuple[str, Path, Dict[str, Path]]]:
        """
        Discover .tf files in the current directory that have corresponding test cases.
        Returns list of tuples: (resource_name, tf_file_path, test_cases_dict)
        """
        resources_with_tests = []

        # Search in the current directory (where tests will be run from)
        current_dir = Path.cwd()

        for tf_file in current_dir.glob("*.tf"):
            # Extract resource name (e.g., "gateways" from "gateways.tf")
            resource_name = tf_file.stem

            # Skip if not in valid resources list (if list is populated)
            if runner.valid_resources and resource_name not in runner.valid_resources:
                print(f"Skipping {tf_file.name} - not in importer.go")
                continue

            # Check if test cases exist for this resource
            test_cases = runner.find_test_cases(resource_name)
            if test_cases:
                resources_with_tests.append((resource_name, tf_file, test_cases))
                print(f"✓ Found test cases for {resource_name}: {list(test_cases.keys())}")
            else:
                # Check if directory exists but files are empty
                resource_test_dir = runner.test_cases_dir / resource_name
                if resource_test_dir.exists():
                    print(f"Test case directory exists for {resource_name} but files are empty or invalid - skipping")
                else:
                    print(f"No test cases found for {resource_name} - skipping")

        return resources_with_tests

    def test_all_resources_lifecycle(self, runner: TerraformTestRunner):
        """
        Test the complete lifecycle of ALL terraform resources in batch mode:
        1. Add ALL test resources from add.tf files (single terraform apply)
        2. Modify ALL test resources from modify.tf files (single terraform apply)
        3. Clean up ALL test resources (single terraform apply)
        """
        print(f"\n{'='*60}")
        print(f"Batch Testing All Resources")
        print(f"{'='*60}")

        # Discover all .tf files with test cases
        resources_with_tests = self.discover_tf_files_with_test_cases(runner)
        
        if not resources_with_tests:
            pytest.skip("No test cases found for any resources")

        print(f"\nDiscovered {len(resources_with_tests)} resources with test cases:")
        for resource_name, tf_file, test_cases in resources_with_tests:
            print(f"  - {resource_name} ({tf_file.name}): {list(test_cases.keys())}")

        # Create backups for all files
        backups = {}
        print(f"\nCreating backups for all files...")
        for resource_name, tf_file, _ in resources_with_tests:
            backup_path = runner.create_backup(tf_file)
            backups[tf_file] = backup_path
        print(f"Created backups in {runner.backup_dir}")

        try:
            # ========================================
            # Pre-flight check: Verify clean state
            # ========================================
            print(f"\n{'='*60}")
            print(f"Pre-flight Check: Verifying clean state")
            print(f"{'='*60}")
            print(f"  Running terraform plan to ensure no pending changes...")
            success, output = runner.terraform_plan(str(runner.tf_dir))
            if not success:
                pytest.fail("terraform plan failed during pre-flight check")
            
            # Check if there are any changes pending
            if "No changes" not in output:
                print(f"\n⚠️  WARNING: Terraform state is not clean!")
                print(f"Output:\n{output}")
                pytest.fail("Terraform has pending changes before test execution. Please ensure clean state.")
            
            print(f"  ✓ Clean state confirmed - no pending changes")

            # ========================================
            # Phase 1: Add ALL test resources
            # ========================================
            add_resources = [(name, tf_file, cases) for name, tf_file, cases in resources_with_tests if "add" in cases]
            
            if add_resources:
                print(f"\n{'='*60}")
                print(f"Phase 1: Adding ALL test resources")
                print(f"{'='*60}")
                
                for resource_name, tf_file, test_cases in add_resources:
                    print(f"  Injecting resources from test_cases/{resource_name}/add.tf → {tf_file.name}")
                    if not runner.inject_resources(tf_file, test_cases["add"]):
                        pytest.fail(f"Failed to inject test resources for {resource_name}")

                # Single terraform apply for all added resources (plan is done automatically)
                print(f"\n  Running terraform apply for all added resources...")
                if not runner.terraform_apply(str(runner.tf_dir)):
                    pytest.fail("terraform apply failed after adding resources")

                # Validate state
                state = runner.terraform_show(str(runner.tf_dir))
                assert state is not None, "Failed to retrieve terraform state"
                print(f"  ✓ All resources added successfully ({len(add_resources)} files processed)")

            # ========================================
            # Phase 2: Modify ALL test resources
            # ========================================
            modify_resources = [(name, tf_file, cases) for name, tf_file, cases in resources_with_tests if "modify" in cases]
            
            if modify_resources:
                print(f"\n{'='*60}")
                print(f"Phase 2: Modifying ALL test resources")
                print(f"{'='*60}")

                # First, remove all previously injected resources
                print(f"  Removing previously injected resources...")
                for resource_name, tf_file, _ in resources_with_tests:
                    runner.remove_injected_resources(tf_file)

                # Now inject modified resources (with merging)
                successfully_merged = []
                for resource_name, tf_file, test_cases in modify_resources:
                    print(f"  Processing modifications for {resource_name}...")
                    
                    # modify.tf MUST have a corresponding add.tf
                    if "add" not in test_cases:
                        print(f"  ⚠️  WARNING: Skipping {resource_name} - modify.tf exists but add.tf is missing!")
                        print(f"      modify.tf can only modify resources that were added via add.tf")
                        continue
                    
                    # Merge modify.tf with add.tf
                    merged_content = runner.merge_modify_with_add(
                        test_cases["add"],
                        test_cases["modify"]
                    )
                    if merged_content:
                        print(f"  Injecting merged resources into {tf_file.name}")
                        try:
                            with open(tf_file, 'a') as f:
                                f.write("\n\n# === TEST RESOURCES START ===\n")
                                f.write(merged_content)
                                f.write("\n# === TEST RESOURCES END ===\n")
                            successfully_merged.append(resource_name)
                        except Exception as e:
                            pytest.fail(f"Failed to write merged resources for {resource_name}: {e}")
                    else:
                        print(f"  ⚠️  WARNING: Skipping {resource_name} - test files are empty or invalid")
                        continue

                # Only run terraform if we successfully merged at least one resource
                if not successfully_merged:
                    print(f"  ⚠️  WARNING: No valid modifications found, skipping Phase 2")
                    print(f"  (All modify.tf files are either empty or have no resources defined)")
                else:
                    # Single terraform apply for all modifications (plan is done automatically)
                    print(f"\n  Running terraform apply for all modifications...")
                    if not runner.terraform_apply(str(runner.tf_dir)):
                        pytest.fail("terraform apply failed after modifying resources")

                    # Validate state
                    state = runner.terraform_show(str(runner.tf_dir))
                    assert state is not None, "Failed to retrieve terraform state"
                    print(f"  ✓ All resources modified successfully ({len(successfully_merged)} files processed)")

            # ========================================
            # Phase 3: Clean up ALL test resources
            # ========================================
            print(f"\n{'='*60}")
            print(f"Phase 3: Cleaning up ALL test resources")
            print(f"{'='*60}")
            
            print(f"  Removing test resources from all files...")
            for resource_name, tf_file, _ in resources_with_tests:
                runner.remove_injected_resources(tf_file)

            # Single terraform apply to remove all resources (plan is done automatically)
            print(f"  Running terraform apply for cleanup...")
            if not runner.terraform_apply(str(runner.tf_dir)):
                pytest.fail("terraform apply failed during cleanup")

            print(f"  ✓ All resources cleaned up successfully")

            # Restore all files from backup
            print(f"\n  Restoring all files from backup...")
            for tf_file, backup_path in backups.items():
                runner.restore_backup(tf_file, backup_path)
            print(f"  ✓ All files restored from backup")

            print(f"\n{'='*60}")
            print(f"✓ All tests completed successfully")
            print(f"  - Total resources tested: {len(resources_with_tests)}")
            print(f"  - Resources with add.tf: {len(add_resources) if add_resources else 0}")
            print(f"  - Resources with modify.tf: {len(modify_resources) if modify_resources else 0}")
            print(f"  - Total terraform applies: 3 (1 per phase)")
            print(f"{'='*60}")

        except Exception as e:
            print(f"\n❌ Test failed: {e}")
            # Restore all backups on failure
            print(f"  Restoring all files from backup after failure...")
            for tf_file, backup_path in backups.items():
                if tf_file.exists() and backup_path.exists():
                    runner.restore_backup(tf_file, backup_path)
            raise
        finally:
            # Ensure we always restore all original files
            print(f"\n  Final cleanup: Ensuring all files are restored...")
            for tf_file, backup_path in backups.items():
                if tf_file.exists() and backup_path.exists():
                    runner.restore_backup(tf_file, backup_path)


if __name__ == "__main__":
    sys.exit(pytest.main([__file__, "-v", "-s"]))
