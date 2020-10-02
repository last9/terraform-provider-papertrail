package papertrail

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oogway/goptrail"
)

func resourcePapertrailSearch() *schema.Resource {
	return &schema.Resource{
		Read:   resourcePapertrailSearchRead,
		Create: resourcePapertrailSearchCreate,
		Update: resourcePapertrailSearchUpdate,
		Delete: resourcePapertrailSearchDelete,

		SchemaVersion: 1,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourcePapertrailSearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(goptrail.Client)

	search, err := client.GetSearch(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", search.Name)
	d.Set("query", search.Query)
	d.Set("group_id", search.Group.ID)
	return nil
}

func resourcePapertrailSearchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(goptrail.Client)

	search := filterSearch(d)

	out, err := client.CreateSearch(search)
	if err != nil {
		return fmt.Errorf("Failed to create search, err: %v", err)
	}

	d.SetId(strconv.Itoa(out.ID))
	d.Set("name", search.Name)
	d.Set("query", search.Query)
	d.Set("group_id", search.Group.ID)

	return nil
}

func resourcePapertrailSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(goptrail.Client)

	search := filterSearch(d)
	search.ID, _ = strconv.Atoi(d.Id())

	return client.UpdateSearch(search)
}

func resourcePapertrailSearchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(goptrail.Client)

	if err := client.DeleteSearch(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func filterSearch(d *schema.ResourceData) goptrail.Search {
	search := goptrail.Search{}
	search.Name = d.Get("name").(string)
	search.Query = d.Get("query").(string)
	search.Group = goptrail.Group{}
	search.Group.ID = d.Get("group_id").(int)

	return search
}
