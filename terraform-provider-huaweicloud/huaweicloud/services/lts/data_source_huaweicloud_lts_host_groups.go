// Generated by PMS #226
package lts

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

func DataSourceLtsHostGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLtsHostGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"host_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Speicifies the ID of the host group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Speicifies the name of the host group.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Speicifies the type of the host group.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the key/value pairs to associate with the host group.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All host groups that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the host group.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the host group.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the host group.`,
						},
						"host_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The ID list of hosts to associate with the host group.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs to associate with the host group.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the host group, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the host group, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

type HostGroupsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newHostGroupsDSWrapper(d *schema.ResourceData, meta interface{}) *HostGroupsDSWrapper {
	return &HostGroupsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceLtsHostGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newHostGroupsDSWrapper(d, meta)
	listHostGroupRst, err := wrapper.ListHostGroup()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listHostGroupToSchema(listHostGroupRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API LTS POST /v3/{project_id}/lts/host-group-list
func (w *HostGroupsDSWrapper) ListHostGroup() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "lts")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/lts/host-group-list"
	params := map[string]any{
		"filter": map[string]any{
			"host_group_name_list": w.PrimToArray("name"),
			"host_group_tag":       w.getFilterHostGroupTag(),
			"host_group_type":      w.Get("type"),
		},
		"host_group_id_list": w.PrimToArray("host_group_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("POST").
		URI(uri).
		Body(params).
		Request().
		Result()
}

func (w *HostGroupsDSWrapper) listHostGroupToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("groups", schemas.SliceToList(body.Get("result"),
			func(groups gjson.Result) any {
				return map[string]any{
					"id":         groups.Get("host_group_id").Value(),
					"name":       groups.Get("host_group_name").Value(),
					"type":       groups.Get("host_group_type").Value(),
					"host_ids":   schemas.SliceToStrList(groups.Get("host_id_list")),
					"tags":       w.setResultHostGroupTag(groups),
					"created_at": w.setResultCreateTime(groups),
					"updated_at": w.setResultUpdateTime(groups),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (w *HostGroupsDSWrapper) getFilterHostGroupTag() map[string]interface{} {
	tags := w.Get("tags")
	if tags != nil {
		return map[string]interface{}{
			"tag_type": "OR",
			"tag_list": utils.ExpandResourceTags(tags.(map[string]interface{})),
		}
	}
	return nil
}

func (*HostGroupsDSWrapper) setResultHostGroupTag(data gjson.Result) map[string]interface{} {
	tags := data.Get("host_group_tag").Value()
	return utils.FlattenTagsToMap(tags)
}

func (*HostGroupsDSWrapper) setResultCreateTime(data gjson.Result) string {
	// Convert to the time corresponding to the local time zone of the computer.
	return utils.FormatTimeStampRFC3339(data.Get("create_time").Int()/1000, false)
}

func (*HostGroupsDSWrapper) setResultUpdateTime(data gjson.Result) string {
	// Convert to the time corresponding to the local time zone of the computer.
	return utils.FormatTimeStampRFC3339(data.Get("update_time").Int()/1000, false)
}
