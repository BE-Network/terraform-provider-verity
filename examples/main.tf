terraform {
  required_providers {
    verity = {
      source = "registry.terraform.io/local/verity"
    }
  }
}

provider "verity" {
  uri      = var.uri
  username = var.username
  password = var.password
}

data "verity_state_importer" "import" {
  output_dir = var.config_dir
}
