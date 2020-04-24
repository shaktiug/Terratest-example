package test

import (
 /* "fmt"
  "io/ioutil"
*/
  "testing"
  "github.com/gruntwork-io/terratest/modules/terraform"
  "github.com/stretchr/testify/assert"
)


func TestTerraformAppService(t *testing.T) {

    terraformOptions := &terraform.Options{
        TerraformDir: "../tf",
}

defer terraform.Destroy(t, terraformOptions)

terraform.InitAndApply(t, terraformOptions)

expected_app_service_default_hostname := "https://beoecomdev-appservice.azurewebsites.net"
expepcted_app_service_name := "beoecomdev-appsevice"

actual_app_service_name := terraform.Output(t, terraformOptions, "app_service_name")
actual_app_service_hostname := terraform.Output(t, terraformOptions, "app_service_default_hostname")

assert.Equal(t, expected_app_service_default_hostname, actual_app_service_hostname)
assert.Equal(t, expepcted_app_service_name, actual_app_service_name)
}
