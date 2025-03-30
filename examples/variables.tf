variable "uri" {
  description = "The base URL of the API"
  type        = string
  sensitive   = true
}

variable "username" {
  description = "API username"
  type        = string
  sensitive   = true
}

variable "password" {
  description = "API password"
  type        = string
  sensitive   = true
}

variable "config_dir" {
  description = "Directory where Terraform configuration files will be generated"
  type        = string
  default     = "."
}