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

func TestAccPapertrailGroup_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPapertrailGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPapertrailGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists("papertrail_group.group"),
					resource.TestCheckResourceAttr("papertrail_group.group", "name", name),
					resource.TestCheckResourceAttr("papertrail_group.group", "system_wildcard", ""),
				),
			},
		},
	})
}

func TestAccPapertrailGroup_with_system_wildcard(t *testing.T) {
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	wildcard := fmt.Sprintf("*/%s/*", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPapertrailGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPapertrailGroupConfigWithWildCard(name, wildcard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists("papertrail_group.group"),
					resource.TestCheckResourceAttr("papertrail_group.group", "name", name),
					resource.TestCheckResourceAttr("papertrail_group.group", "system_wildcard", wildcard),
				),
			},
		},
	})
}

func testAccCheckGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Group ID is set")
		}

		conn := testAccProvider.Meta().(goptrail.Client)
		group, err := conn.GetGroup(rs.Primary.ID)
		if err != nil {
			return err
		}

		if group.ID != rs.Primary.ID {
			return fmt.Errorf("Incorrect Group ID: %d", group.ID)
		}
		return nil
	}
}

func testAccCheckPapertrailGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(goptrail.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "papertrail_group" {
			continue
		}

		group, err := conn.GetGroup(rs.Primary.ID)
		if !strings.Contains(err.Error(), ":Not Found") {
			return err
		}

		if group.ID != "" {
			return fmt.Errorf("Group exists, ID: %d", group.ID)
		}
	}
	return nil
}

func testAccPapertrailGroupConfig(name string) string {
	return fmt.Sprintf(`resource "papertrail_group" "group" {
  name             = "%s"
}`, name)
}

func testAccPapertrailGroupConfigWithWildCard(name, wc string) string {
	return fmt.Sprintf(`resource "papertrail_group" "group" {
  name             = "%s"
  system_wildcard  = "%s"
}`, name, wc)
}
