package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAssignmentPackages_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_assignment_packages.basic"
	dataSource2 := "data.huaweicloud_rms_assignment_packages.filter_by_name"
	dataSource3 := "data.huaweicloud_rms_assignment_packages.filter_by_status"
	dataSource4 := "data.huaweicloud_rms_assignment_packages.filter_by_id"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsAssignmentPackages_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsAssignmentPackages_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_assignment_packages" "basic" {
  depends_on = [huaweicloud_rms_assignment_package.test]
}

data "huaweicloud_rms_assignment_packages" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_rms_assignment_package.test]
}

data "huaweicloud_rms_assignment_packages" "filter_by_status" {
  status = "CREATE_SUCCESSFUL"

  depends_on = [huaweicloud_rms_assignment_package.test]
}

data "huaweicloud_rms_assignment_packages" "filter_by_id" {
  package_id = huaweicloud_rms_assignment_package.test.id
}

locals {
  name_filter_result   = [for v in data.huaweicloud_rms_assignment_packages.filter_by_name.packages[*].name : v == "%[2]s"]
  status_filter_result = [for v in data.huaweicloud_rms_assignment_packages.filter_by_name.packages[*].status : v == "CREATE_SUCCESSFUL"]
  id_filter_result     = [
    for v in data.huaweicloud_rms_assignment_packages.filter_by_name.packages[*].id : v == huaweicloud_rms_assignment_package.test.id
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_assignment_packages.basic.packages) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}
`, testAssignmentPackage_basic(name), name)
}
