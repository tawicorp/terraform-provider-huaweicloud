package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/cdn/v1/domains"
)

func TestAccCdnDomain_basic(t *testing.T) {
	var domain domains.CdnDomain

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCDN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainV1Exists("huaweicloud_cdn_domain_v1.domain_1", &domain),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain_v1.domain_1", "name", HW_CDN_DOMAIN_NAME),
				),
			},
		},
	})
}

func testAccCheckCdnDomainV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cdnClient, err := config.CdnV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CDN Domain client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cdn_domain_v1" {
			continue
		}

		found, err := domains.Get(cdnClient, rs.Primary.ID, nil).Extract()
		if err == nil && found.DomainStatus != "deleting" {
			return fmt.Errorf("Destroying CDN domain failed or domain still exists")
		}
	}

	return nil
}

func testAccCheckCdnDomainV1Exists(n string, domain *domains.CdnDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("CDN Domain Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cdnClient, err := config.CdnV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CDN Domain client: %s", err)
		}

		found, err := domains.Get(cdnClient, rs.Primary.ID, nil).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("CDN Domain not found")
		}

		*domain = *found
		return nil
	}
}

var testAccCdnDomainV1_basic = fmt.Sprintf(`
resource "huaweicloud_cdn_domain_v1" "domain_1" {
  name   = "%s"
  type   = "web"
  enterprise_project_id = 0
  service_area = "outside_mainland_china"
  sources {
      active = 1
      origin = "100.254.53.75"
      origin_type  = "ipaddr"
  }
}
`, HW_CDN_DOMAIN_NAME)
