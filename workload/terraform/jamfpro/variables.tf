variable "username" {
  description = "JAMF Pro api username for Terraform automation"
  type        = string
}

variable "password" {
  description = "JAMF Pro api password for Terraform automation"
  type        = string
}

variable "server_url" {
  description = "This is the JAMF Pro instance for terraform automation. Value should be in the format xxxx.jamfcloud.com"
  type        = string
}