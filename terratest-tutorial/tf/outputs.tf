output "app_service_name" {
   value = "${azurerm_app_service.appsvc.name}"
}

output "app_service_default_hostname" { 
  value = "https://${azurerm_app_service.appsvc.default_site_hostname}"
}
