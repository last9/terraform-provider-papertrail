package papertrail

import (
	"fmt"
	"strconv"

	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oogway/goptrail"
)

func dataSourcePapertrailLogDestination() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourcePapertrailLogDestinationRead,
		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePapertrailLogDestinationRead(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(goptrail.Client)
	if !ok {
		return fmt.Errorf("Cannot convert %v to PapertrailClient", meta)
	}

	var id, port int

	if val, ok := d.GetOk("id"); ok {
		id = val.(int)
	}

	if val, ok := d.GetOk("port"); ok {
		port = val.(int)
	}

	if id == 0 && port == 0 {
		return errors.New("one of id or port is required")
	}

	destinations, err := client.ListLogDestinations()
	if err != nil {
		return err
	}

	for _, dest := range destinations {
		if dest.ID == id || dest.Syslog.Port == port {
			d.SetId(strconv.Itoa(dest.ID))
			d.Set("hostname", dest.Syslog.Hostname)
			d.Set("port", dest.Syslog.Port)
			d.Set("description", dest.Syslog.Description)
			return nil
		}
	}

	return fmt.Errorf("Log Destination with id %v or port %v found in %v", id, port, destinations)
}
