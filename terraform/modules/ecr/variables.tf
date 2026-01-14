variable "image_registry" {
  type = string
}

variable "image_mutable" {
  type = string
  default = "MUTABLE"
}

variable "image_scan_on_push" {
  type = bool
  default = true
}