package test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const grant_type string = "client_credentials"
const client_id string = "9f4c2aac-9db0-444b-90cb-988d0138e4dc"
const client_secret string = "O9bTKN:HQ6b=Zq?xnr.PCCq50wx3L/ev"
const resource string = "https://management.azure.com"
const redirect_uri string = "https://ShaktiServicePrincipal"
const tenant_id string = "29425ebd-9122-4d2b-aeaa-1257aec8b162"
const subscription_id string = "08896c78-cdb8-472b-bef3-b474be3e57fc"
const api_endpoint string = "https://login.microsoftonline.com/29425ebd-9122-4d2b-aeaa-1257aec8b162/oauth2/token"
const appid string = client_id
const rg string = "beoecomtest"
const name string = "beoecomdev-appservice"

func get_token() string {
	formData := url.Values{
		"grant_type":    {grant_type},
		"client_id":     {client_id},
		"client_secret": {client_secret},
		"redirect_uri":  {redirect_uri},
		"resource":      {resource},
	}
	resp, err := http.PostForm(api_endpoint, formData)
	if err != nil {
		fmt.Printf("Error for http/PostForm() %s\n", err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["access_token"].(string)
}

func list_app_service_pub_profile(tkn string) string {
	verb := "POST"
	app_service_pub_sig := "https://management.azure.com/subscriptions/%s/resourceGroups/" + "%s/providers/Microsoft.Web/sites/%s/publishxml?api-version=2019-08-01"

	app_service_pub_url := fmt.Sprintf(app_service_pub_sig, subscription_id, rg, name)

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + tkn

	// Create a new request using http
	req, err := http.NewRequest(verb, app_service_pub_url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return (string([]byte(body)))
}

func get_publish_app_svc_endpoints(tkn string) (string, string, string, string) {

	pub_profile_xml := list_app_service_pub_profile(tkn)

	type PublishData struct {
		XMLName        xml.Name `xml:"publishData"`
		Text           string   `xml:",chardata"`
		PublishProfile []struct {
			Text                        string `xml:",chardata"`
			ProfileName                 string `xml:"profileName,attr"`
			PublishMethod               string `xml:"publishMethod,attr"`
			PublishUrl                  string `xml:"publishUrl,attr"`
			MsdeploySite                string `xml:"MsdeploySite,attr"`
			UserName                    string `xml:"userName,attr"`
			UserPwd                     string `xml:"userPWD,attr"`
			DestinationAppUrl           string `xml:"destinationAppUrl,attr"`
			SQLServerDBConnectionString string `xml:"SQLServerDBConnectionString,attr"`
			MySQLDBConnectionString     string `xml:"mySQLDBConnectionString,attr"`
			HostingProviderForumLink    string `xml:"hostingProviderForumLink,attr"`
			ControlPanelLink            string `xml:"controlPanelLink,attr"`
			WebSystem                   string `xml:"webSystem,attr"`
			FtpPassiveMode              string `xml:"ftpPassiveMode,attr"`
			Databases                   string `xml:"databases"`
		} `xml:"publishProfile"`
	}

	var publishData PublishData
	err := xml.Unmarshal([]byte(pub_profile_xml), &publishData)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "", "", "", ""
	}

	var appUrl, publishUrl, userName, userPwd string = "", "", "", ""

	for _, elem := range publishData.PublishProfile {

		if !strings.HasSuffix(elem.ProfileName, "ReadOnly - FTP") &&
			strings.HasSuffix(elem.ProfileName, "- FTP") {

			publishUrl = elem.PublishUrl
			userName = elem.UserName
			userPwd = elem.UserPwd
			appUrl = elem.DestinationAppUrl

			break

		}
	}
	return appUrl, publishUrl, userName, userPwd
}

func upload_local_file(webpage string, publishUrl string, username string, password string) {

	name_pwd := fmt.Sprintf("%s:%s", username, password)
	url := fmt.Sprintf("%s/", publishUrl)
	curl := exec.Command("curl", "-s", "-T", webpage, "-u", name_pwd, url)
	curl.Stdout = os.Stdout
	curl.Stderr = os.Stderr
	err := curl.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func test_app_service(appurl, webpage string) bool {
	webpageurl := fmt.Sprintf("%s/%s", appurl, webpage)
	verb := "GET"

	// Create a new request using http
	req, err := http.NewRequest(verb, webpageurl, nil)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	s := string([]byte(body))
	return strings.Index(s, "Shakti's Azure WebApp") != -1
}

func TestTerraformAppService(t *testing.T) {

	terraformOptions := &terraform.Options{
		TerraformDir: "../tf",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	expected_app_service_default_hostname := "https://beoecomdev-appservice.azurewebsites.net"
	expepcted_app_service_name := "beoecomdev-appservice"

	actual_app_service_name := terraform.Output(t, terraformOptions, "app_service_name")
	actual_app_service_hostname := terraform.Output(t, terraformOptions, "app_service_default_hostname")

	assert.Equal(t, expected_app_service_default_hostname, actual_app_service_hostname)
	assert.Equal(t, expepcted_app_service_name, actual_app_service_name)

	fmt.Println("(1) Getting Auth Token")
	ad_token := get_token()

	fmt.Println("(2) Getting PublishUrl, Username, password for FTP upload to App Service")

	appUrl, publishUrl, userName, password := get_publish_app_svc_endpoints(ad_token)
	fmt.Println("(3) Uploading index.html")
	var webpage string = "index.html"

	upload_local_file(webpage, publishUrl, userName, password)
	fmt.Println("(4) Testing index.html")

	//Make sure app is up and running
	assert.Equal(t, test_app_service(appUrl, webpage), true)

}
