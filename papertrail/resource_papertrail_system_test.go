package papertrail

import (
	"errors"
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/oogway/goptrail"
)

func TestAccPapertrailSystem_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	hostname := fmt.Sprintf("%s.hostname.com", acctest.RandString(10))
	destination_port := 514

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPapertrailSystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPapertrailSystemConfig(name, hostname, destination_port),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSystemExists("papertrail_system.system"),
					resource.TestCheckResourceAttr("papertrail_system.system", "name", name),
					resource.TestCheckResourceAttr("papertrail_system.system", "hostname", hostname),
				),
			},
		},
	})
}

func testAccCheckSystemExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No System ID is set")
		}

		conn := testAccProvider.Meta().(goptrail.Client)
		system, err := conn.GetSystem(rs.Primary.ID)
		if err != nil {
			return err
		}

		if system.ID != rs.Primary.ID {
			return fmt.Errorf("Incorrect System ID: %d", system.ID)
		}
		return nil
	}
}

func testAccCheckPapertrailSystemDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(goptrail.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "papertrail_system" {
			continue
		}

		system, err := conn.GetSystem(rs.Primary.ID)
		if !strings.Contains(err.Error(), ":Not Found") {
			return err
		}

		if system.ID != "" {
			return fmt.Errorf("System exists, ID: %d", system.ID)
		}
	}
	return nil
}

func testAccPapertrailSystemConfig(name, hostname string, destination_port int) string {
	return fmt.Sprintf(`resource "papertrail_system" "system" {
  name             = "%s"
  hostname         = "%s"
  destination_port = %d
}`, name, hostname, destination_port)
}
