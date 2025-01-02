// Generated by PMS #359
package antiddos

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

func DataSourceWeeklyProtectionStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWeeklyProtectionStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"period_start_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start date of the seven-day period, the value is a timestamp in milliseconds.`,
			},
			"ddos_intercept_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of DDoS attacks blocked in a week.`,
			},
			"weekdata": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The number of attacks in a week.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period_start_date": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The start date of the seven-day period, the value is a timestamp in milliseconds.`,
						},
						"ddos_intercept_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of DDoS attacks blocked.`,
						},
						"ddos_blackhole_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of DDoS black holes.`,
						},
						"max_attack_bps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum attack traffic.`,
						},
						"max_attack_conns": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of attack connections.`,
						},
					},
				},
			},
			"top10": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Top 10 attacked IP addresses.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"floating_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Elastic IP address.`,
						},
						"times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of DDoS attacks blocked, including scrubbing and black holes.`,
						},
					},
				},
			},
		},
	}
}

type WeeklyProtectionStatisticsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newWeeklyProtectionStatisticsDSWrapper(d *schema.ResourceData, meta interface{}) *WeeklyProtectionStatisticsDSWrapper {
	return &WeeklyProtectionStatisticsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceWeeklyProtectionStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newWeeklyProtectionStatisticsDSWrapper(d, meta)
	listWeeklyReportsRst, err := wrapper.ListWeeklyReports()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listWeeklyReportsToSchema(listWeeklyReportsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API ANTI-DDOS GET /v1/{project_id}/antiddos/weekly
func (w *WeeklyProtectionStatisticsDSWrapper) ListWeeklyReports() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "anti-ddos")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/antiddos/weekly"
	params := map[string]any{
		"period_start_date": w.Get("period_start_date"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *WeeklyProtectionStatisticsDSWrapper) listWeeklyReportsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("ddos_intercept_times", body.Get("ddos_intercept_times").Value()),
		d.Set("weekdata", schemas.SliceToList(body.Get("weekdata"),
			func(weekdata gjson.Result) any {
				return map[string]any{
					"period_start_date":    weekdata.Get("period_start_date").Value(),
					"ddos_intercept_times": weekdata.Get("ddos_intercept_times").Value(),
					"ddos_blackhole_times": weekdata.Get("ddos_blackhole_times").Value(),
					"max_attack_bps":       weekdata.Get("max_attack_bps").Value(),
					"max_attack_conns":     weekdata.Get("max_attack_conns").Value(),
				}
			},
		)),
		d.Set("top10", schemas.SliceToList(body.Get("top10"),
			func(top10 gjson.Result) any {
				return map[string]any{
					"floating_ip_address": top10.Get("floating_ip_address").Value(),
					"times":               top10.Get("times").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
