// Generated by PMS #194
package as

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceAsNotifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAsNotificationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the AS group to which the notifications belong.`,
			},
			"topic_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the topic name in SMN.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All AS group notifications that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The topic name in SMN.`,
						},
						"topic_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique topic URN in SMN.`,
						},
						"events": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The notification scene list.`,
						},
					},
				},
			},
		},
	}
}

type NotificationsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newNotificationsDSWrapper(d *schema.ResourceData, meta interface{}) *NotificationsDSWrapper {
	return &NotificationsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceAsNotificationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newNotificationsDSWrapper(d, meta)
	lisScaNotRst, err := wrapper.ListScalingNotifications()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listScalingNotificationsToSchema(lisScaNotRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}
func (w *NotificationsDSWrapper) ListScalingNotifications() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "autoscaling")
	if err != nil {
		return nil, err
	}

	uri := "/autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}"
	uri = strings.ReplaceAll(uri, "{scaling_group_id}", w.Get("scaling_group_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Filter(
			filters.New().From("topics").
				Where("topic_name", "=", w.Get("topic_name")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *NotificationsDSWrapper) listScalingNotificationsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("topics", schemas.SliceToList(body.Get("topics"),
			func(topics gjson.Result) any {
				return map[string]any{
					"topic_name": topics.Get("topic_name").Value(),
					"topic_urn":  topics.Get("topic_urn").Value(),
					"events":     schemas.SliceToStrList(topics.Get("topic_scene")),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}