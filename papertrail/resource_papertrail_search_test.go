package papertrail

import (
	"errors"
	"fmt"
	"testing"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oogway/goptrail"
)

func TestAccPapertrailSearch_basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	query := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPapertrailSearchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPapertrailSearchConfig(name, query),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSearchExists("papertrail_search.search"),
					resource.TestCheckResourceAttr("papertrail_search.search", "name", name),
					resource.TestCheckResourceAttr("papertrail_search.search", "query", query),
				),
			},
		},
	})
}

func testAccCheckSearchExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Search ID is set")
		}

		conn := testAccProvider.Meta().(goptrail.Client)
		search, err := conn.GetSearch(rs.Primary.ID)
		if err != nil {
			return err
		}

		if strconv.Itoa(search.ID) != rs.Primary.ID {
			return fmt.Errorf("Incorrect Search ID: %d", search.ID)
		}
		return nil
	}
}

func testAccCheckPapertrailSearchDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(goptrail.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "papertrail_search" {
			continue
		}

		_, err := conn.GetSearch(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Search still exists: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccPapertrailSearchConfig(name, query string) string {
	return fmt.Sprintf(`
resource "papertrail_group" "group" {
	name             = "tf-acc-test-group"
	system_wildcard  = "*"
}

resource "papertrail_search" "search" {
  name     = "%s"
  query    = "%s"
  group_id = "${papertrail_group.group.id}"
}
`, name, query)
}
