package papertrail

import (
	"errors"
	"fmt"
	"testing"

	"strconv"

	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/oogway/goptrail"
)

func TestAccPapertrailSystemGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPapertrailSystemGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSystemGroupExists("papertrail_system_group.psg"),
					resource.TestCheckResourceAttrSet("papertrail_system_group.psg", "system_id"),
					resource.TestCheckResourceAttrSet("papertrail_system_group.psg", "group_id"),
				),
			},
		},
	})
}

func testAccCheckSystemGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Group ID is set")
		}

		client := testAccProvider.Meta().(goptrail.Client)
		group, err := client.GetGroup(rs.Primary.Attributes["group_id"])
		if err != nil {
			return err
		}

		var found bool
		for _, sys := range group.Systems {
			if strconv.Itoa(sys.ID) == rs.Primary.Attributes["system_id"] {
				found = true
				break
			}
		}

		if !found {
			return errors.New("System not found in Group")
		}
		return nil
	}
}

func testAccCheckPapertrailSystemGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(goptrail.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "papertrail_system_group" {
			continue
		}

		group, err := client.GetGroup(rs.Primary.Attributes["group_id"])
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil
		} else if err != nil {
			return err
		}

		for _, sys := range group.Systems {
			if strconv.Itoa(sys.ID) == rs.Primary.Attributes["system_id"] {
				return errors.New("system group is not deleted")
			}
		}
	}

	return nil
}

func testAccPapertrailSystemGroupConfig() string {
	return fmt.Sprintf(`resource "papertrail_system" "system" {
  name             = "%s"
  destination_port = 29504
}

resource "papertrail_group" "group" {
  name             = "%s"
  system_wildcard  = "%s"
}

resource "papertrail_system_group" "psg" {
  system_id        = "${papertrail_system.system.id}"
  group_id         = "${papertrail_group.group.id}"
}`, acctest.RandString(4), acctest.RandString(4), acctest.RandString(4))
}
