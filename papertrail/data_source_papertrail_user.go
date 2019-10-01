package papertrail

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oogway/goptrail"
)

func dataSourcePapertrailUser() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourcePapertrailUserRead,
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePapertrailUserRead(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	email, ok := d.Get("email").(string)
	if !ok {
		return fmt.Errorf("Cannot convert email %v to string", d.Get("email"))
	}

	users, err := client.ListUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Email == email {
			d.SetId(strconv.Itoa(u.ID))
			return nil
		}
	}

	return fmt.Errorf("User with email %v not found in %v", email, users)
}
