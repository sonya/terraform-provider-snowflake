package datasources

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var roleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The role for which to return metadata.",
	},
	"comment": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The comment on the role",
	},
}

// Role Snowflake Role resource.
func Role() *schema.Resource {
	return &schema.Resource{
		Read:               ReadRole,
		Schema:             roleSchema,
		DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_roles instead.",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadRole Reads the database metadata information.
func ReadRole(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	roleName := d.Get("name").(string)

	role, err := client.Roles.ShowByID(ctx, sdk.NewAccountObjectIdentifier(roleName))
	if err != nil {
		log.Printf("[DEBUG] role (%s) not found", roleName)
		d.SetId("")
		return nil
	}

	d.SetId(role.Name)
	if err := d.Set("name", role.Name); err != nil {
		return err
	}
	if err := d.Set("comment", role.Comment); err != nil {
		return err
	}
	return nil
}
