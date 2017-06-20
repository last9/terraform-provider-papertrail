package papertrail

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oogway/goptrail"
)

func resourcePapertrailSystem() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePapertrailSystemRead,
		Create: resourcePapertrailSystemCreate,
		Update: resourcePapertrailSystemUpdate,
		Delete: resourcePapertrailSystemDelete,

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
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"destination_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_event_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"syslog_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"syslog_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePapertrailSystemRead(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	system, err := client.GetSystem(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", system.Name)
	d.Set("ip_address", system.IpAddress)
	d.Set("hostname", system.Hostname)
	d.Set("last_event_at", system.LastEventAt)
	d.Set("syslog_hostname", system.Syslog.Hostname)
	d.Set("syslog_port", system.Syslog.Port)

	return nil
}

func resourcePapertrailSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	system, err := filterSystem(d)
	if err != nil {
		return err
	}

	// in case some resources has changed
	system.ID, _ = strconv.Atoi(d.Id())

	needUpdate := false
	for _, val := range []string{
		"name",
		"ip_address",
		"hostname",
		"destination_id",
		"destination_port",
	} {
		if d.HasChange(val) {
			needUpdate = true
		}
	}

	if needUpdate {
		if err := client.UpdateSystem(system); err != nil {
			return err
		}
	}
	return resourcePapertrailSystemRead(d, meta)
}

func resourcePapertrailSystemCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	system, err := filterSystem(d)
	if err != nil {
		return err
	}

	out, err := client.RegisterSystem(system)
	if err != nil {
		return fmt.Errorf("Failed to register system: %+v", err)
	}

	d.SetId(strconv.Itoa(out.ID))

	return resourcePapertrailSystemRead(d, meta)
}

func resourcePapertrailSystemDelete(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	if err := client.UnregisterSystem(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func filterSystem(d *schema.ResourceData) (goptrail.InputSystem, error) {
	system := goptrail.InputSystem{}
	system.Name = d.Get("name").(string)

	if val, ok := d.GetOk("ip_address"); ok {
		system.Hostname = val.(string)
	}

	if val, ok := d.GetOk("hostname"); ok {
		system.Hostname = val.(string)
	}

	if val, ok := d.GetOk("destination_id"); ok {
		system.DestinationID = val.(int)
	}

	if val, ok := d.GetOk("destination_port"); ok {
		system.DestinationPort = val.(int)
	}

	if val, ok := d.GetOk("description"); ok {
		system.Description = val.(string)
	}

	if system.DestinationID == 0 && system.DestinationPort == 0 {
		return system, fmt.Errorf("Either the destination_id or destination_port must be specified")
	}

	return system, nil
}
