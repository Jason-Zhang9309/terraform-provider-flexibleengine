package flexibleengine

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/addons"
)

func TestAccCCEAddonV3_basic(t *testing.T) {
	var addon addons.Addon

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "flexibleengine_cce_addon_v3.test"
	clusterName := "flexibleengine_cce_cluster_v3.cluster_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEAddonV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonV3Exists(resourceName, clusterName, &addon),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
		},
	})
}

func testAccCheckCCEAddonV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceAddonV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating flexibleengine CCE Addon client: %s", err)
	}

	var clusterId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "flexibleengine_cce_cluster_v3" {
			clusterId = rs.Primary.ID
		}

		if rs.Type != "flexibleengine_cce_addon_v3" {
			continue
		}

		if clusterId != "" {
			_, err := addons.Get(cceClient, rs.Primary.ID, clusterId).Extract()
			if err == nil {
				return fmt.Errorf("addon still exists")
			}
		}
	}
	return nil
}

func testAccCheckCCEAddonV3Exists(n string, cluster string, addon *addons.Addon) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmt.Errorf("Cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		if c.Primary.ID == "" {
			return fmt.Errorf("Cluster id is not set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceAddonV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating flexibleengine CCE Addon client: %s", err)
		}

		found, err := addons.Get(cceClient, rs.Primary.ID, c.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Addon not found")
		}

		*addon = *found

		return nil
	}
}

func testAccCCEAddonV3_basic(rName string) string {
	return fmt.Sprintf(`
resource "flexibleengine_cce_cluster_v3" "cluster_1" {
	name = "%s"
	cluster_type="VirtualMachine"
	flavor_id="cce.s1.small"
	vpc_id="%s"
	subnet_id="%s"
	container_network_type="overlay_l2"
	}
	
resource "flexibleengine_cce_node_v3" "node_1" {
	cluster_id = "${flexibleengine_cce_cluster_v3.cluster_1.id}"
	name = "test-node"
	flavor_id="s1.medium"
	availability_zone= "%s"
	key_pair="%s"
	root_volume {
		size= 40
		volumetype= "SATA"
	}
	data_volumes {
		size= 100
		volumetype= "SATA"
	}
	}

resource "flexibleengine_cce_addon_v3" "test" {
    cluster_id = flexibleengine_cce_cluster_v3.cluster_1.id
    version = "1.0.3"
	template_name = "metrics-server"
	depends_on = [flexibleengine_cce_node_v3.node_1]
}
`, rName, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE, OS_KEYPAIR_NAME)
}
