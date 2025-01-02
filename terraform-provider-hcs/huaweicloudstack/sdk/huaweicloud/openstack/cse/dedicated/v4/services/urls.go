package services

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("registry", "microservices")
}

func resourceURL(client *golangsdk.ServiceClient, serviceId string) string {
	return client.ServiceURL("registry", "microservices", serviceId)
}
