from test_terraform_resources import TerraformTestRunner


def pytest_addoption(parser):
    """Register custom --mode command-line option."""
    parser.addoption(
        "--mode",
        required=True,
        choices=TerraformTestRunner.VALID_MODES,
        help="Provider mode: datacenter or campus",
    )
