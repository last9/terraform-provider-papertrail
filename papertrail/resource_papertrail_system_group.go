package papertrail

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oogway/goptrail"
)

func resourcePapertrailSystemGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePapertrailSystemGroupRead,
		Create: resourcePapertrailSystemGroupCreate,
		Update: resourcePapertrailSystemGroupUpdate,
		Delete: resourcePapertrailSystemGroupDelete,

		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"system_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePapertrailSystemGroupRead(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	systemID := d.Get("system_id").(string)
	groupID := d.Get("group_id").(string)

	group, err := client.GetGroup(groupID)
	if err != nil {
		return err
	}

	for ix, sys := range group.Systems {
		if strconv.Itoa(sys.ID) == systemID {
			d.Set("index", ix)
			return nil
		}
	}

	return fmt.Errorf("System: %s is not found in Group: %s", systemID, groupID)
}

func resourcePapertrailSystemGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	systemID := d.Get("system_id").(string)
	groupID := d.Get("group_id").(string)

	if err := client.AddSystemToGroup(systemID, groupID); err != nil {
		return fmt.Errorf("Failed to create group, err: %v", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", systemID, groupID))

	return resourcePapertrailSystemGroupRead(d, meta)
}

func resourcePapertrailSystemGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	var osd, nsd, ogd, ngd string
	var sIDChanged, gIDChanged bool

	nsd = d.Get("system_id").(string)
	ngd = d.Get("group_id").(string)

	if d.HasChange("system_id") {
		sIDChanged = true
		oc, nc := d.GetChange("system_id")
		osd = oc.(string)
		nsd = nc.(string)
	}

	if d.HasChange("group_id") {
		gIDChanged = true
		oc, nc := d.GetChange("group_id")
		ogd = oc.(string)
		ngd = nc.(string)
	}

	if sIDChanged && gIDChanged {
		if err := client.RemoveSystemFromGroup(osd, ogd); err != nil {
			return fmt.Errorf("Failed to remove system: %s from group: %s, err: %v",
				osd, ogd, err)
		}

		if err := client.AddSystemToGroup(nsd, ngd); err != nil {
			return fmt.Errorf("Failed to add system: %s to group: %s, err: %v",
				nsd, ngd, err)
		}

		return nil
	}

	if sIDChanged {
		if err := client.RemoveSystemFromGroup(osd, ngd); err != nil {
			return fmt.Errorf("Failed to remove system: %s from group: %s, err: %v",
				osd, ngd, err)
		}

		if err := client.AddSystemToGroup(nsd, ngd); err != nil {
			return fmt.Errorf("Failed to add system: %s to group: %s, err: %v",
				nsd, ngd, err)
		}

		return nil
	}

	if gIDChanged {
		if err := client.RemoveSystemFromGroup(nsd, ogd); err != nil {
			return fmt.Errorf("Failed to remove system: %s from group: %s, err: %v",
				nsd, ogd, err)
		}

		if err := client.AddSystemToGroup(nsd, ngd); err != nil {
			return fmt.Errorf("Failed to add system: %s to group: %s, err: %v",
				nsd, ngd, err)
		}

		return nil
	}

	return nil
}

func resourcePapertrailSystemGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	systemID := d.Get("system_id").(string)
	groupID := d.Get("group_id").(string)

	if err := client.RemoveSystemFromGroup(systemID, groupID); err != nil {
		return fmt.Errorf("Failed to remove system: %s from group: %s, err: %v",
			systemID, groupID, err)
	}

	d.SetId("")

	return nil
}
