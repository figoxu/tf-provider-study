// Generated by PMS #361
package dbss

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
)

func DataSourceDbssAuditDataMaskingRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbssAuditDataMaskingRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the audit instance ID to which the privacy data masking rules belong.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the privacy data masking rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the privacy data masking rule.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the privacy data masking rule.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the privacy data masking rule.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the privacy data masking rule.`,
						},
						"regex": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The regular expression of the privacy data masking rule.`,
						},
						"mask_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The privacy data display substitution value.`,
						},
						"operate_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The operation time of the privacy data masking rule, in UTC format.`,
						},
					},
				},
			},
		},
	}
}

type AuditDataMaskingRulesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAuditDataMaskingRulesDSWrapper(d *schema.ResourceData, meta interface{}) *AuditDataMaskingRulesDSWrapper {
	return &AuditDataMaskingRulesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDbssAuditDataMaskingRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAuditDataMaskingRulesDSWrapper(d, meta)
	lisAudSenMasRst, err := wrapper.ListAuditSensitiveMasks()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAuditSensitiveMasksToSchema(lisAudSenMasRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DBSS GET /v1/{project_id}/{instance_id}/dbss/audit/sensitive/masks
func (w *AuditDataMaskingRulesDSWrapper) ListAuditSensitiveMasks() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dbss")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/{instance_id}/dbss/audit/sensitive/masks"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("rules", "offset", "limit", 100).
		Request().
		Result()
}

func (w *AuditDataMaskingRulesDSWrapper) listAuditSensitiveMasksToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("rules", schemas.SliceToList(body.Get("rules"),
			func(rules gjson.Result) any {
				return map[string]any{
					"id":           rules.Get("id").Value(),
					"name":         rules.Get("name").Value(),
					"type":         rules.Get("type").Value(),
					"status":       rules.Get("status").Value(),
					"regex":        rules.Get("regex").Value(),
					"mask_value":   rules.Get("mask_value").Value(),
					"operate_time": rules.Get("operate_time").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
