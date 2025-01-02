// Generated by PMS #443
package dbss

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDbssRdsDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbssRdsDatabasesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the RDS database type.`,
			},
			"databases": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The RDS database list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS instance ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS database name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS database type.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS instance status.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS database version.`,
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS database IP address.`,
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS database port.`,
						},
						"is_supported": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether agent-free audit is supported.`,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The RDS instance name.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID to which the RDS instance belongs.`,
						},
					},
				},
			},
		},
	}
}

type RdsDatabasesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRdsDatabasesDSWrapper(d *schema.ResourceData, meta interface{}) *RdsDatabasesDSWrapper {
	return &RdsDatabasesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDbssRdsDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRdsDatabasesDSWrapper(d, meta)
	listRdsDatabasesRst, err := wrapper.ListRdsDatabases()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listRdsDatabasesToSchema(listRdsDatabasesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DBSS GET /v2/{project_id}/audit/databases/rds
func (w *RdsDatabasesDSWrapper) ListRdsDatabases() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dbss")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/audit/databases/rds"
	params := map[string]any{
		"db_type": w.Get("type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("databases", "offset", "limit", 0).
		Request().
		Result()
}

func (w *RdsDatabasesDSWrapper) listRdsDatabasesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("databases", schemas.SliceToList(body.Get("databases"),
			func(databases gjson.Result) any {
				return map[string]any{
					"id":                    databases.Get("id").Value(),
					"name":                  databases.Get("db_name").Value(),
					"type":                  databases.Get("type").Value(),
					"status":                databases.Get("status").Value(),
					"version":               databases.Get("version").Value(),
					"ip":                    databases.Get("ip").Value(),
					"port":                  databases.Get("port").Value(),
					"is_supported":          databases.Get("is_supported").Value(),
					"instance_name":         databases.Get("instance_name").Value(),
					"enterprise_project_id": databases.Get("enterprise_id").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
