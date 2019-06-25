package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAwsServiceQuotasServiceQuota_basic(t *testing.T) {
	resourceName := "aws_servicequotas_service_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSServiceQuotas(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsServiceQuotasServiceQuotaConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_code", "L-F678F1CE"),
					resource.TestCheckResourceAttr(resourceName, "service_code", "vpc"),
					resource.TestCheckResourceAttr(resourceName, "value", "75"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsServiceQuotasServiceQuota_Value_IncreaseOnCreate(t *testing.T) {
	quotaCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_CREATE_QUOTA_CODE")
	if quotaCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_CREATE_QUOTA_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	serviceCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_CREATE_SERVICE_CODE")
	if serviceCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_CREATE_SERVICE_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	value := os.Getenv("SERVICEQUOTAS_INCREASE_ON_CREATE_VALUE")
	if value == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_CREATE_VALUE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	resourceName := "aws_servicequotas_service_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSServiceQuotas(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsServiceQuotasServiceQuotaConfigValue(quotaCode, serviceCode, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_code", quotaCode),
					resource.TestCheckResourceAttr(resourceName, "service_code", serviceCode),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
			},
		},
	})
}

func TestAccAwsServiceQuotasServiceQuota_Value_IncreaseOnUpdate(t *testing.T) {
	t.Skip("Requires aws_servicequotas_service_quota data source")

	quotaCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_UPDATE_QUOTA_CODE")
	if quotaCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_UPDATE_QUOTA_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	serviceCode := os.Getenv("SERVICEQUOTAS_INCREASE_ON_UPDATE_SERVICE_CODE")
	if serviceCode == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_UPDATE_SERVICE_CODE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	value := os.Getenv("SERVICEQUOTAS_INCREASE_ON_UPDATE_VALUE")
	if value == "" {
		t.Skip(
			"Environment variable SERVICEQUOTAS_INCREASE_ON_UPDATE_VALUE is not set. " +
				"WARNING: This test will submit a real service quota increase!")
	}

	dataSourceName := "aws_servicequotas_service_quota.test"
	resourceName := "aws_servicequotas_service_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSServiceQuotas(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsServiceQuotasServiceQuotaConfigSameValue(quotaCode, serviceCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_code", quotaCode),
					resource.TestCheckResourceAttr(resourceName, "service_code", serviceCode),
					resource.TestCheckResourceAttrPair(resourceName, "value", dataSourceName, "value"),
				),
			},
			{
				Config: testAccAwsServiceQuotasServiceQuotaConfigValue(quotaCode, serviceCode, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "quota_code", quotaCode),
					resource.TestCheckResourceAttr(resourceName, "service_code", serviceCode),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
			},
		},
	})
}

func testAccPreCheckAWSServiceQuotas(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).servicequotasconn

	input := &servicequotas.ListServicesInput{}

	_, err := conn.ListServices(input)

	if testAccPreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccAwsServiceQuotasServiceQuotaConfig() string {
	return fmt.Sprintf(`
resource "aws_servicequotas_service_quota" "test" {
  quota_code   = "L-F678F1CE"
  service_code = "vpc"
  value        = 75
}
`)
}

func testAccAwsServiceQuotasServiceQuotaConfigSameValue(quotaCode, serviceCode string) string {
	return fmt.Sprintf(`
data "aws_servicequotas_service_quota" "test" {
  quota_code   = %[1]q
  service_code = %[2]q
}

resource "aws_servicequotas_service_quota" "test" {
  quota_code   = "${data.aws_servicequotas_service_quota.test.quota_code}"
  service_code = "${data.aws_servicequotas_service_quota.test.service_code}"
  value        = "${data.aws_servicequotas_service_quota.test.value}"
}
`, quotaCode, serviceCode)
}

func testAccAwsServiceQuotasServiceQuotaConfigValue(quotaCode, serviceCode, value string) string {
	return fmt.Sprintf(`
resource "aws_servicequotas_service_quota" "test" {
  quota_code   = %[1]q
  service_code = %[2]q
  value        = %[3]s
}
`, quotaCode, serviceCode, value)
}
