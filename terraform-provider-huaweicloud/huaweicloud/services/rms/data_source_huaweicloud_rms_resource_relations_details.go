// Generated by PMS #438
package rms

import (
	"context"
	"strings"

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

func DataSourceRmsRelationsDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRmsRelationsDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource ID.`,
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the direction of a resource relationship.`,
			},
			"relations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resource relationships.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The relationship type.`,
						},
						"from_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the source resource.`,
						},
						"to_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the destination resource.`,
						},
						"from_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source resource ID.`,
						},
						"to_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The destination resource ID.`,
						},
					},
				},
			},
		},
	}
}

type RelationsDetailsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRelationsDetailsDSWrapper(d *schema.ResourceData, meta interface{}) *RelationsDetailsDSWrapper {
	return &RelationsDetailsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRmsRelationsDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRelationsDetailsDSWrapper(d, meta)
	shoResRelDetRst, err := wrapper.ShowResourceRelationsDetail()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showResourceRelationsDetailToSchema(shoResRelDetRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/all-resources/{resource_id}/relations
func (w *RelationsDetailsDSWrapper) ShowResourceRelationsDetail() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rms")
	if err != nil {
		return nil, err
	}

	uri := "/v1/resource-manager/domains/{domain_id}/all-resources/{resource_id}/relations"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Config.DomainID)
	uri = strings.ReplaceAll(uri, "{resource_id}", w.Get("resource_id").(string))
	params := map[string]any{
		"direction": w.Get("direction"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("relations", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *RelationsDetailsDSWrapper) showResourceRelationsDetailToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("relations", schemas.SliceToList(body.Get("relations"),
			func(relations gjson.Result) any {
				return map[string]any{
					"relation_type":      relations.Get("relation_type").Value(),
					"from_resource_type": relations.Get("from_resource_type").Value(),
					"to_resource_type":   relations.Get("to_resource_type").Value(),
					"from_resource_id":   relations.Get("from_resource_id").Value(),
					"to_resource_id":     relations.Get("to_resource_id").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}