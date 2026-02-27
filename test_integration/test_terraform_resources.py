#!/usr/bin/env python3
import os
import sys
import shutil
import subprocess
import json
import tempfile
import re
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Tuple, Optional, Set, Any

try:
    import pytest
except ImportError:
    print("ERROR: pytest is not installed")
    sys.exit(1)

try:
    import hcl2
except ImportError:
    print("ERROR: python-hcl2 is not installed")
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
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        self.api_log_path = self.tf_dir / f"test_logs_{timestamp}.log"
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

    def run_terraform_command(
        self, command: List[str], cwd: str,
        env: Optional[Dict[str, str]] = None
    ) -> Tuple[int, str, str]:
        """Run a terraform command and return exit code, stdout, stderr."""
        try:
            result = subprocess.run(
                command,
                cwd=cwd,
                capture_output=True,
                text=True,
                timeout=300,
                env=env
            )
            return result.returncode, result.stdout, result.stderr
        except subprocess.TimeoutExpired:
            return -1, "", "Command timed out after 300 seconds"
        except Exception as e:
            return -1, "", str(e)

    def terraform_plan(self, cwd: Optional[str] = None) -> Tuple[bool, str]:
        """Run terraform plan and return success status and output."""
        cwd = cwd or str(self.tf_dir)
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

    def terraform_apply(
        self, cwd: Optional[str] = None,
        phase_label: Optional[str] = None
    ) -> bool:
        """Apply terraform changes and capture API request logs."""
        cwd = cwd or str(self.tf_dir)
        print(f"  Running terraform apply in {cwd}...")

        # Set up debug logging to capture API requests
        tmp = tempfile.NamedTemporaryFile(
            suffix='.log', prefix='tf_debug_', delete=False
        )
        debug_log = Path(tmp.name)
        tmp.close()

        env = os.environ.copy()
        env['TF_LOG'] = 'DEBUG'
        env['TF_LOG_PATH'] = str(debug_log)

        returncode, stdout, stderr = self.run_terraform_command(
            ["terraform", "apply", "-auto-approve", "tfplan"],
            cwd,
            env=env
        )

        # Extract and log API requests from debug output
        if debug_log.exists():
            exchanges = self._extract_api_exchanges(debug_log)
            if exchanges:
                label = phase_label or "terraform apply"
                self._log_api_exchanges(exchanges, label)
                print(
                    f"  ✓ Logged {len(exchanges)} API request(s)"
                    f" to {self.api_log_path.name}"
                )
            try:
                debug_log.unlink()
            except Exception:
                pass

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

    def terraform_plan_and_apply(
        self, cwd: Optional[str] = None,
        phase_label: Optional[str] = None
    ) -> bool:
        """Run terraform plan followed by terraform apply. Returns True on success."""
        cwd = cwd or str(self.tf_dir)
        success, output = self.terraform_plan(cwd)
        if not success:
            return False
        return self.terraform_apply(cwd, phase_label=phase_label)

    def terraform_show(self, cwd: Optional[str] = None) -> Optional[Dict]:
        """Get current terraform state as JSON."""
        cwd = cwd or str(self.tf_dir)
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

    def inject_content(self, target_file: Path, content: str) -> bool:
        """Inject content string into the target file between markers."""
        try:
            with open(target_file, 'a') as f:
                f.write("\n\n# === TEST RESOURCES START ===\n")
                f.write(content)
                f.write("\n# === TEST RESOURCES END ===\n")
            return True
        except Exception as e:
            print(f"  ❌ Failed to inject content: {e}")
            return False

    def inject_resources(self, target_file: Path, test_case_file: Path) -> bool:
        """Inject test resources from a file into the target file."""
        try:
            with open(test_case_file, 'r') as f:
                test_content = f.read()
            return self.inject_content(target_file, test_content)
        except Exception as e:
            print(f"  ❌ Failed to read test case file: {e}")
            return False

    def remove_injected_resources(self, target_file: Path) -> bool:
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

    @staticmethod
    def _is_indexed_block_list(value: Any) -> bool:
        """Check if a value is a list of dicts where each dict has an 'index' key."""
        return (
            isinstance(value, list)
            and len(value) > 0
            and all(isinstance(item, dict) and 'index' in item for item in value)
        )

    def _merge_indexed_blocks(self, base_list: List[Dict], override_list: List[Dict]) -> List[Dict]:
        """
        Merge two lists of indexed block dicts by matching on the 'index' field.
        - Matched indexes: deep-merge attributes (override wins)
        - Unmatched base indexes: preserved as-is
        - New override indexes: appended
        Result is sorted by index.
        """
        # Build lookup by index
        merged = {item['index']: item.copy() for item in base_list}

        for override_item in override_list:
            idx = override_item['index']
            if idx in merged:
                # Deep-merge the matched block
                merged[idx] = self._merge_hcl_attributes(merged[idx], override_item)
            else:
                # New block from override
                merged[idx] = override_item.copy()

        # Return sorted by index
        return [merged[k] for k in sorted(merged.keys())]

    def _merge_hcl_attributes(self, base: Dict[str, Any], override: Dict[str, Any]) -> Dict[str, Any]:
        """
        Merge HCL attributes, with override taking precedence.
        - dict + dict: recursive deep merge
        - list-of-indexed-dicts + list-of-indexed-dicts: merge by index field
        - anything else: full replacement
        """
        merged = base.copy()

        for key, value in override.items():
            if key in merged and isinstance(merged[key], dict) and isinstance(value, dict):
                # Recursively merge nested dictionaries
                merged[key] = self._merge_hcl_attributes(merged[key], value)
            elif (key in merged
                  and self._is_indexed_block_list(merged[key])
                  and self._is_indexed_block_list(value)):
                # Index-aware merge for nested blocks with 'index' field
                merged[key] = self._merge_indexed_blocks(merged[key], value)
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

    def _find_resource_block_span(
        self, content: str, resource_type: str, resource_name: str
    ) -> Optional[Tuple[int, int]]:
        """
        Find the start and end positions of a resource block in HCL text.
        Returns (start, end) character positions, or None if not found.
        """
        pattern = rf'resource\s+"{re.escape(resource_type)}"\s+"{re.escape(resource_name)}"\s*\{{'
        match = re.search(pattern, content)
        if not match:
            return None

        # Walk forward from the opening brace counting braces to find the matching close
        brace_start = match.end() - 1  # Position of the opening {
        depth = 1
        pos = brace_start + 1
        while pos < len(content) and depth > 0:
            ch = content[pos]
            if ch == '{':
                depth += 1
            elif ch == '}':
                depth -= 1
            elif ch == '"':
                # Skip string contents to avoid counting braces inside strings
                pos += 1
                while pos < len(content) and content[pos] != '"':
                    if content[pos] == '\\':
                        pos += 1  # Skip escaped character
                    pos += 1
            elif ch == '#':
                # Skip single-line comments
                while pos < len(content) and content[pos] != '\n':
                    pos += 1
            elif ch == '/' and pos + 1 < len(content) and content[pos + 1] == '/':
                # Skip // line comments
                while pos < len(content) and content[pos] != '\n':
                    pos += 1
            elif ch == '/' and pos + 1 < len(content) and content[pos + 1] == '*':
                # Skip /* */ block comments
                pos += 2
                while pos + 1 < len(content) and not (content[pos] == '*' and content[pos + 1] == '/'):
                    pos += 1
                pos += 1  # Skip past the closing /
            pos += 1

        if depth != 0:
            return None

        return (match.start(), pos)

    def _replace_resource_block_in_text(
        self, content: str, resource_type: str, resource_name: str, new_hcl: str
    ) -> Optional[str]:
        """
        Replace a resource block in HCL text with new content.
        Returns modified text, or None if the resource block was not found.
        """
        span = self._find_resource_block_span(content, resource_type, resource_name)
        if span is None:
            return None
        start, end = span
        return content[:start] + new_hcl + content[end:]

    def merge_modify_with_existing_tf(self, tf_file: Path, modify_file: Path) -> bool:
        """
        For update-only resources: merge modify.tf changes into the existing .tf file.
        Parses the existing resource from the .tf file as the base, merges with modify.tf
        overrides, and replaces the resource block in-place.
        Returns True on success, False on failure.
        """
        try:
            # Parse modify.tf to get override resources
            with open(modify_file, 'r') as f:
                modify_content = hcl2.load(f)

            if 'resource' not in modify_content or not modify_content['resource']:
                print(f"  ⚠️  WARNING: modify.tf has no resources defined")
                return False

            # Read the existing .tf file text (for block replacement)
            with open(tf_file, 'r') as f:
                tf_text = f.read()

            # Parse the existing .tf file to get base resources
            with open(tf_file, 'r') as f:
                tf_parsed = hcl2.load(f)

            # Build lookup of existing resources
            existing_resources = {}
            if 'resource' in tf_parsed:
                for resource_list in tf_parsed['resource']:
                    for resource_type, resources in resource_list.items():
                        for resource_name, attributes in resources.items():
                            key = (resource_type, resource_name)
                            existing_resources[key] = attributes

            # Process each resource in modify.tf
            modified_tf_text = tf_text
            modified_count = 0
            for resource_list in modify_content['resource']:
                for resource_type, resources in resource_list.items():
                    for resource_name, attributes in resources.items():
                        key = (resource_type, resource_name)
                        if key not in existing_resources:
                            print(
                                f"  ⚠️  WARNING: Resource {resource_type}.{resource_name}"
                                f" from modify.tf not found in {tf_file.name} - skipping"
                            )
                            continue

                        # Merge: existing base + modify overrides
                        merged_attrs = self._merge_hcl_attributes(
                            existing_resources[key],
                            attributes
                        )

                        # Format as HCL
                        new_hcl = self._format_hcl_resource(resource_type, resource_name, merged_attrs)

                        # Replace the resource block in-place
                        result = self._replace_resource_block_in_text(
                            modified_tf_text, resource_type, resource_name, new_hcl
                        )
                        if result is None:
                            print(
                                f"  ❌ Failed to locate resource block"
                                f" {resource_type}.{resource_name} in {tf_file.name}"
                            )
                            return False
                        modified_tf_text = result
                        modified_count += 1
                        print(f"  ✓ Merged modifications for {resource_type}.{resource_name}")

            if modified_count == 0:
                print(f"  ⚠️  WARNING: No resources were modified in {tf_file.name}")
                return False

            # Write back the modified file
            with open(tf_file, 'w') as f:
                f.write(modified_tf_text)

            return True

        except Exception as e:
            print(f"  ❌ Failed to merge modify.tf with {tf_file.name}: {e}")
            return False

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
                                # New resource in modify.tf — not in add.tf
                                print(
                                    f"  ⚠️  WARNING: {resource_type}.{resource_name}"
                                    f" exists in modify.tf but not in add.tf — adding as-is"
                                )
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
        """
        Find test case files for a given resource and validate they contain resources.
        Returns a dict with optional keys: 'add', 'modify'.
        If only 'modify' is present (no valid 'add'), the resource is update-only.
        """
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

    @staticmethod
    def is_update_only(test_cases: Dict[str, Path]) -> bool:
        """Check if a resource is update-only (has modify.tf but no add.tf)."""
        return "modify" in test_cases and "add" not in test_cases

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

    def discover_resources(self, verbose: bool = False) -> List[Tuple[str, Path, Dict[str, Path]]]:
        """
        Discover .tf files in the current directory that have corresponding test cases.
        Returns list of tuples: (resource_name, tf_file_path, test_cases_dict)
        """
        resources_with_tests = []
        current_dir = Path.cwd()

        for tf_file in sorted(current_dir.glob("*.tf")):
            resource_name = tf_file.stem

            if self.valid_resources and resource_name not in self.valid_resources:
                if verbose:
                    print(f"Skipping {tf_file.name} - not in importer.go")
                continue

            test_cases = self.find_test_cases(resource_name)
            if test_cases:
                resources_with_tests.append((resource_name, tf_file, test_cases))
                if verbose:
                    label = "[update-only] " if self.is_update_only(test_cases) else ""
                    print(f"✓ Found {label}test case for {resource_name}: {list(test_cases.keys())}")
            else:
                if verbose:
                    resource_test_dir = self.test_cases_dir / resource_name
                    if resource_test_dir.exists():
                        print(
                            f"Test case directory exists for {resource_name}"
                            f" but files are empty or invalid - skipping"
                        )
                    else:
                        print(f"No test cases found for {resource_name} - skipping")

        return resources_with_tests

    def start_api_log(self, test_name: str):
        """Start a fresh API request log for this test run."""
        with open(self.api_log_path, 'w') as f:
            f.write(f"{'='*80}\n")
            f.write(f"API Request Log: {test_name}\n")
            f.write(f"{'='*80}\n")
        print(f"  API request log: {self.api_log_path}")

    def _extract_api_exchanges(
        self, log_path: Path
    ) -> List[Dict[str, str]]:
        """Extract PUT/PATCH/DELETE HTTP exchanges from terraform debug log."""
        try:
            content = log_path.read_text(errors='replace')
        except Exception:
            return []

        exchanges: List[Dict[str, str]] = []
        current: Optional[Dict[str, str]] = None

        prefix_re = re.compile(
            r'^\d{4}-\d{2}-\d{2}T(\d{2}:\d{2}:\d{2})\.\d+[+-]\d+'
            r' \[DEBUG\] provider\.\S+:\s*(.*)'
        )
        request_re = re.compile(
            r'^(PUT|PATCH|DELETE) (/\S+) HTTP/'
        )
        response_re = re.compile(r'^HTTP/\d+\.\d+ (\d+ .+)')

        for line in content.splitlines():
            m = prefix_re.match(line)
            if not m:
                continue

            timestamp, body = m.group(1), m.group(2)

            # New PUT/PATCH/DELETE request?
            req_m = request_re.match(body)
            if req_m:
                # Save any previous incomplete exchange
                if current:
                    exchanges.append(current)
                current = {
                    'timestamp': timestamp,
                    'method': req_m.group(1),
                    'url': req_m.group(2),
                    'request_body': '',
                    'response_status': '',
                    'response_body': '',
                    '_phase': 'request',
                }
                continue

            if not current:
                continue

            # Response status line?
            resp_m = response_re.match(body)
            if resp_m:
                current['response_status'] = resp_m.group(1)
                current['_phase'] = 'response'
                continue

            # JSON body line?
            stripped = body.strip()
            if stripped.startswith('{'):
                if (current['_phase'] == 'request'
                        and not current['request_body']):
                    current['request_body'] = stripped
                elif (current['_phase'] == 'response'
                        and not current['response_body']):
                    current['response_body'] = stripped
                    # Exchange complete
                    exchanges.append(current)
                    current = None

        # Capture any trailing incomplete exchange
        if current:
            exchanges.append(current)

        return exchanges

    def _log_api_exchanges(
        self, exchanges: List[Dict[str, str]],
        phase_label: str
    ):
        """Append extracted API exchanges to the log file."""
        with open(self.api_log_path, 'a') as f:
            f.write(f"\n{'─'*80}\n")
            f.write(f"{phase_label}\n")
            f.write(f"{'─'*80}\n\n")

            for ex in exchanges:
                f.write(
                    f"[{ex['timestamp']}] "
                    f"{ex['method']} {ex['url']}\n"
                )
                if ex['request_body']:
                    f.write(f"  → {ex['request_body']}\n")
                if ex['response_status']:
                    f.write(f"  ← {ex['response_status']}\n")
                if ex['response_body']:
                    f.write(f"  ← {ex['response_body']}\n")
                f.write("\n")

    def cleanup_plan_file(self):
        """Remove leftover tfplan file if it exists."""
        tfplan_path = self.tf_dir / "tfplan"
        if tfplan_path.exists():
            try:
                tfplan_path.unlink()
            except Exception:
                pass

    def cleanup(self):
        """Clean up backup directory."""
        if self.backup_dir and self.backup_dir.exists():
            shutil.rmtree(self.backup_dir)


def _discover_resources() -> List[Tuple[str, str, Dict[str, str]]]:
    """
    Discover .tf files in the current directory that have corresponding test cases.
    Called at pytest collection time for parametrization of single-resource tests.
    Returns list of tuples with string-serialized paths for clean pytest IDs.
    """
    runner = TerraformTestRunner()
    resources = []
    for name, tf_file, test_cases in runner.discover_resources(verbose=False):
        cases_str = {k: str(v) for k, v in test_cases.items()}
        resources.append((name, str(tf_file), cases_str))
    runner.cleanup()
    return resources


def pytest_generate_tests(metafunc):
    """Dynamically parametrize test_single_resource with discovered resources."""
    if "resource_info" in metafunc.fixturenames:
        resources = _discover_resources()
        metafunc.parametrize(
            "resource_info",
            resources,
            ids=[name for name, _, _ in resources],
        )


def pytest_collection_modifyitems(config, items):
    """
    Auto-deselect single-resource tests unless the user passed -k.
    This keeps 'test_all_resources_lifecycle' as the default test,
    while 'test_single_resource[X]' only runs when explicitly selected.
    """
    keyword_expr = config.getoption("-k", default="")
    if keyword_expr:
        return

    # No -k flag: deselect all single_resource tests so only the batch test runs
    remaining = []
    deselected = []
    for item in items:
        if "test_single_resource" in item.nodeid:
            deselected.append(item)
        else:
            remaining.append(item)
    if deselected:
        config.hook.pytest_deselected(items=deselected)
        items[:] = remaining


class TestTerraformResources:
    """Pytest test class for terraform resources."""

    @pytest.fixture(scope="class")
    def runner(self):
        """Create a test runner instance."""
        runner = TerraformTestRunner()
        yield runner
        runner.cleanup()

    # ==================================================================
    # DEFAULT: Batch test — all resources in a single run
    # ==================================================================

    def test_all_resources_lifecycle(self, runner: TerraformTestRunner):
        """
        Test the complete lifecycle of ALL terraform resources in batch mode:
        - Normal resources: Add (create) → Modify (update) → Cleanup (delete)
        - Update-only resources: Modify (update existing) → Revert (restore original)
        All phases use single terraform apply calls for efficiency.
        """
        print(f"\n{'='*60}")
        print(f"Batch Testing All Resources")
        print(f"{'='*60}")

        resources_with_tests = runner.discover_resources(verbose=True)

        if not resources_with_tests:
            pytest.skip("No test cases found for any resources")

        runner.start_api_log("Batch Test: All Resources")

        normal_resources = [(n, f, c) for n, f, c in resources_with_tests if not runner.is_update_only(c)]
        update_only_resources = [(n, f, c) for n, f, c in resources_with_tests if runner.is_update_only(c)]

        print(f"\nDiscovered {len(resources_with_tests)} resources with test cases:")
        for resource_name, tf_file, test_cases in normal_resources:
            print(f"  - {resource_name} ({tf_file.name}): {list(test_cases.keys())}")
        for resource_name, tf_file, test_cases in update_only_resources:
            print(f"  - {resource_name} ({tf_file.name}): {list(test_cases.keys())} [update-only]")

        # Create backups for all files
        backups = {}
        print(f"\nCreating backups for all files...")
        for resource_name, tf_file, _ in resources_with_tests:
            backup_path = runner.create_backup(tf_file)
            backups[tf_file] = backup_path
        print(f"Created backups in {runner.backup_dir}")

        try:
            print(f"\n{'='*60}")
            print(f"Pre-flight Check: Verifying clean state")
            print(f"{'='*60}")
            print(f"  Running terraform plan to ensure no pending changes...")
            success, output = runner.terraform_plan()
            if not success:
                pytest.fail("terraform plan failed during pre-flight check")

            if "No changes" not in output:
                print(f"\n⚠️  WARNING: Terraform state is not clean!")
                print(f"Output:\n{output}")
                pytest.fail("Terraform has pending changes before test execution. Please ensure clean state.")

            runner.cleanup_plan_file()
            print(f"  ✓ Clean state confirmed - no pending changes")

            # Phase 1: Add ALL test resources (normal resources only)
            add_resources = [(name, tf_file, cases) for name, tf_file, cases in normal_resources if "add" in cases]

            if add_resources:
                print(f"\n{'='*60}")
                print(f"Phase 1: Adding ALL test resources")
                print(f"{'='*60}")

                for resource_name, tf_file, test_cases in add_resources:
                    print(f"  Injecting resources from test_cases/{resource_name}/add.tf → {tf_file.name}")
                    if not runner.inject_resources(tf_file, test_cases["add"]):
                        pytest.fail(f"Failed to inject test resources for {resource_name}")

                print(f"\n  Running terraform plan+apply for all added resources...")
                if not runner.terraform_plan_and_apply(
                    phase_label="Phase 1: Add"
                ):
                    pytest.fail("terraform plan+apply failed after adding resources")

                state = runner.terraform_show()
                assert state is not None, "Failed to retrieve terraform state"
                print(f"  ✓ All resources added successfully ({len(add_resources)} files processed)")
            else:
                print(f"\n  No normal resources with add.tf found, skipping Phase 1")

            # Phase 2: Modify ALL test resources
            modify_normal = [(n, f, c) for n, f, c in normal_resources if "modify" in c]

            has_modifications = modify_normal or update_only_resources
            if has_modifications:
                print(f"\n{'='*60}")
                print(f"Phase 2: Modifying ALL test resources")
                print(f"{'='*60}")

                # Handle normal resources
                successfully_merged = []
                if modify_normal:
                    print(f"  Removing previously injected resources...")
                    for resource_name, tf_file, _ in normal_resources:
                        runner.remove_injected_resources(tf_file)

                    for resource_name, tf_file, test_cases in modify_normal:
                        print(f"  Processing modifications for {resource_name}...")

                        if "add" not in test_cases:
                            print(f"  ⚠️  WARNING: Skipping {resource_name} - modify.tf exists but add.tf is missing!")
                            continue

                        merged_content = runner.merge_modify_with_add(
                            test_cases["add"],
                            test_cases["modify"]
                        )
                        if merged_content:
                            print(f"  Injecting merged resources into {tf_file.name}")
                            if not runner.inject_content(tf_file, merged_content):
                                pytest.fail(f"Failed to write merged resources for {resource_name}")
                            successfully_merged.append(resource_name)
                        else:
                            print(f"  ⚠️  WARNING: Skipping {resource_name} - test files are empty or invalid")
                            continue

                # Handle update-only resources
                successfully_updated = []
                if update_only_resources:
                    print(f"\n  Processing update-only resources...")
                    for resource_name, tf_file, test_cases in update_only_resources:
                        print(f"  Processing update-only modifications for {resource_name}...")
                        if runner.merge_modify_with_existing_tf(tf_file, test_cases["modify"]):
                            successfully_updated.append(resource_name)
                        else:
                            print(f"  ⚠️  WARNING: Skipping update-only resource {resource_name} - merge failed")

                if not successfully_merged and not successfully_updated:
                    print(f"  ⚠️  WARNING: No valid modifications found, skipping Phase 2")
                else:
                    print(f"\n  Running terraform plan+apply for all modifications...")
                    if not runner.terraform_plan_and_apply(
                        phase_label="Phase 2: Modify"
                    ):
                        pytest.fail("terraform plan+apply failed after modifying resources")

                    state = runner.terraform_show()
                    assert state is not None, "Failed to retrieve terraform state"
                    total_modified = len(successfully_merged) + len(successfully_updated)
                    print(f"  ✓ All resources modified successfully ({total_modified} files processed)")
                    if successfully_updated:
                        print(
                            f"    (including {len(successfully_updated)} update-only:"
                            f" {', '.join(successfully_updated)})"
                        )

            # Phase 3: Clean up ALL test resources + revert update-only
            print(f"\n{'='*60}")
            print(f"Phase 3: Cleaning up ALL test resources")
            print(f"{'='*60}")

            print(f"  Removing test resources from normal resource files...")
            for resource_name, tf_file, _ in normal_resources:
                runner.remove_injected_resources(tf_file)

            if update_only_resources:
                print(f"  Reverting update-only resources to original state...")
                for resource_name, tf_file, _ in update_only_resources:
                    if tf_file in backups:
                        runner.restore_backup(tf_file, backups[tf_file])
                        print(f"    ✓ Restored {tf_file.name} from backup")

            print(f"  Running terraform plan+apply for cleanup...")
            if not runner.terraform_plan_and_apply(
                phase_label="Phase 3: Cleanup"
            ):
                pytest.fail("terraform plan+apply failed during cleanup")

            print(f"  ✓ All resources cleaned up successfully")

            print(f"\n{'='*60}")
            print(f"✓ All tests completed successfully")
            print(f"  - Total resources tested: {len(resources_with_tests)}")
            print(f"  - Normal resources (add→modify→delete): {len(normal_resources)}")
            print(f"  - Update-only resources (modify→revert): {len(update_only_resources)}")
            print(f"{'='*60}")

        except Exception as e:
            print(f"\n❌ Test failed: {e}")
            raise
        finally:
            # Always restore all original files regardless of success or failure
            print(f"\n  Final cleanup: Restoring all files from backup...")
            for tf_file, backup_path in backups.items():
                if tf_file.exists() and backup_path.exists():
                    runner.restore_backup(tf_file, backup_path)
            print(f"  ✓ All files restored")

    # ==================================================================
    # SINGLE-RESOURCE: Parametrized test — runs only with -k flag
    # ==================================================================

    def test_single_resource(self, runner: TerraformTestRunner, resource_info):
        """
        Test a single resource's lifecycle in isolation.
        Auto-deselected unless -k is used.

        Usage:
            pytest -k badges -v -s    # test only badges
            pytest -k sites -v -s     # test only sites
        """
        resource_name, tf_file_str, test_cases_str = resource_info
        tf_file = Path(tf_file_str)
        test_cases = {k: Path(v) for k, v in test_cases_str.items()}
        is_update_only = runner.is_update_only(test_cases)

        label = f"{resource_name} [update-only]" if is_update_only else resource_name
        print(f"\n{'='*60}")
        print(f"Single-resource test: {label}")
        print(f"  tf_file: {tf_file.name}")
        print(f"  test_cases: {list(test_cases.keys())}")
        print(f"{'='*60}")

        runner.start_api_log(f"Single Resource: {label}")
        backup_path = runner.create_backup(tf_file)

        try:
            print(f"\n  Pre-flight: Verifying clean state...")
            success, output = runner.terraform_plan()
            if not success:
                pytest.fail("terraform plan failed during pre-flight check")
            if "No changes" not in output:
                print(f"\n⚠️  WARNING: Terraform state is not clean!")
                print(f"Output:\n{output}")
                pytest.fail("Terraform has pending changes before test execution. Please ensure clean state.")

            runner.cleanup_plan_file()
            print(f"  ✓ Clean state confirmed")

            if is_update_only:
                self._run_single_update_only(runner, resource_name, tf_file, test_cases, backup_path)
            else:
                self._run_single_normal(runner, resource_name, tf_file, test_cases)

            print(f"\n{'='*60}")
            print(f"✓ {label} — test completed successfully")
            print(f"{'='*60}")

        except Exception as e:
            print(f"\n❌ Test failed for {resource_name}: {e}")
            raise
        finally:
            if tf_file.exists() and backup_path.exists():
                runner.restore_backup(tf_file, backup_path)

    @staticmethod
    def _run_single_normal(
        runner: TerraformTestRunner,
        resource_name: str,
        tf_file: Path,
        test_cases: Dict[str, Path],
    ):
        """Add → Modify → Cleanup for one normal resource."""

        # Phase 1: Add
        if "add" in test_cases:
            print(f"\n  --- Phase 1: Add ---")
            print(f"  Injecting resources from test_cases/{resource_name}/add.tf → {tf_file.name}")
            if not runner.inject_resources(tf_file, test_cases["add"]):
                pytest.fail(f"Failed to inject test resources for {resource_name}")

            print(f"  Running terraform plan+apply...")
            if not runner.terraform_plan_and_apply(
                phase_label=f"Phase 1: Add [{resource_name}]"
            ):
                pytest.fail("terraform plan+apply failed after adding resources")

            state = runner.terraform_show()
            assert state is not None, "Failed to retrieve terraform state"
            print(f"  ✓ Resources added successfully")

        # Phase 2: Modify
        if "modify" in test_cases and "add" in test_cases:
            print(f"\n  --- Phase 2: Modify ---")
            runner.remove_injected_resources(tf_file)

            merged_content = runner.merge_modify_with_add(
                test_cases["add"],
                test_cases["modify"],
            )
            if merged_content:
                print(f"  Injecting merged resources into {tf_file.name}")
                if not runner.inject_content(tf_file, merged_content):
                    pytest.fail(f"Failed to write merged resources for {resource_name}")

                print(f"  Running terraform plan+apply...")
                if not runner.terraform_plan_and_apply(
                    phase_label=f"Phase 2: Modify [{resource_name}]"
                ):
                    pytest.fail("terraform plan+apply failed after modifying resources")

                state = runner.terraform_show()
                assert state is not None, "Failed to retrieve terraform state"
                print(f"  ✓ Resources modified successfully")
            else:
                print(f"  ⚠️  modify.tf is empty or invalid — skipping Phase 2")

        # Phase 3: Cleanup
        print(f"\n  --- Phase 3: Cleanup ---")
        runner.remove_injected_resources(tf_file)

        print(f"  Running terraform plan+apply for cleanup...")
        if not runner.terraform_plan_and_apply(
            phase_label=f"Phase 3: Cleanup [{resource_name}]"
        ):
            pytest.fail("terraform plan+apply failed during cleanup")
        print(f"  ✓ Resources cleaned up successfully")

    @staticmethod
    def _run_single_update_only(
        runner: TerraformTestRunner,
        resource_name: str,
        tf_file: Path,
        test_cases: Dict[str, Path],
        backup_path: Path,
    ):
        """Modify → Revert for one update-only resource."""

        # Phase 1: Modify existing resource in-place
        print(f"\n  --- Phase 1: Modify (update-only) ---")
        if not runner.merge_modify_with_existing_tf(tf_file, test_cases["modify"]):
            pytest.fail(f"Failed to merge modify.tf into {tf_file.name}")

        print(f"  Running terraform plan+apply...")
        if not runner.terraform_plan_and_apply(
            phase_label=f"Phase 1: Modify [{resource_name}] (update-only)"
        ):
            pytest.fail("terraform plan+apply failed after modifying update-only resource")

        state = runner.terraform_show()
        assert state is not None, "Failed to retrieve terraform state"
        print(f"  ✓ Update-only resource modified successfully")

        # Phase 2: Revert to original
        print(f"\n  --- Phase 2: Revert (update-only) ---")
        runner.restore_backup(tf_file, backup_path)
        print(f"  Restored {tf_file.name} from backup")

        print(f"  Running terraform plan+apply to revert...")
        if not runner.terraform_plan_and_apply(
            phase_label=f"Phase 2: Revert [{resource_name}] (update-only)"
        ):
            pytest.fail("terraform plan+apply failed during revert of update-only resource")
        print(f"  ✓ Update-only resource reverted successfully")


if __name__ == "__main__":
    sys.exit(pytest.main([__file__, "-v", "-s"]))
