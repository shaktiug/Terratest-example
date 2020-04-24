provider "azurerm" {
   version = "=2.5.0"
   features {}
}

resource "azurerm_app_service_plan" "asp" {
  name = "beoecomtest-asp"
  location = "East US"
  resource_group_name = "beoecomtest"
  kind = "Linux"
  reserved = "true"

  sku {
     tier = "Standard"
     size = "S1"
 }
}

resource "azurerm_app_service" "appsvc" {
  name = "beoecomdev-appservice"
  location = "East US"
  resource_group_name = "beoecomtest"
  app_service_plan_id = "${azurerm_app_service_plan.asp.id}"
  }

