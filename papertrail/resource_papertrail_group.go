package papertrail

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oogway/goptrail"
)

func resourcePapertrailGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePapertrailGroupRead,
		Create: resourcePapertrailGroupCreate,
		Update: resourcePapertrailGroupUpdate,
		Delete: resourcePapertrailGroupDelete,

		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_wildcard": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePapertrailGroupRead(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	group, err := client.GetGroup(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", group.Name)
	d.Set("system_wildcard", group.SystemWildcard)
	return nil
}

func resourcePapertrailGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	group := filterGroup(d)

	out, err := client.CreateGroup(group)
	if err != nil {
		return fmt.Errorf("Failed to create group, err: %v", err)
	}

	d.SetId(strconv.Itoa(out.ID))

	return resourcePapertrailGroupRead(d, meta)
}

func resourcePapertrailGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	group := filterGroup(d)

	group.ID, _ = strconv.Atoi(d.Id())

	if d.HasChange("name") || d.HasChange("system_wildcard") {
		if err := client.UpdateGroup(group); err != nil {
			return err
		}
	}

	return resourcePapertrailGroupRead(d, meta)
}

func resourcePapertrailGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	if err := client.DeleteGroup(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func filterGroup(d *schema.ResourceData) goptrail.Group {
	group := goptrail.Group{}
	group.Name = d.Get("name").(string)

	if val, ok := d.GetOk("system_wildcard"); ok {
		group.SystemWildcard = val.(string)
	}
	return group
}
