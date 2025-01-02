// Generated by PMS #198
package dms

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

func DataSourceDmsRocketmqTopicConsumerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketmqTopicConsumerGroupsRead,

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
				Description: `Specifies the instance ID.`,
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the topic name.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the consumer group list.`,
			},
		},
	}
}

type RocketmqTopicConsumerGroupsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRocketmqTopicConsumerGroupsDSWrapper(d *schema.ResourceData, meta interface{}) *RocketmqTopicConsumerGroupsDSWrapper {
	return &RocketmqTopicConsumerGroupsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDmsRocketmqTopicConsumerGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRocketmqTopicConsumerGroupsDSWrapper(d, meta)
	lisConGroOfTopRst, err := wrapper.ListConsumerGroupOfTopic()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listConsumerGroupOfTopicToSchema(lisConGroOfTopRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/topics/{topic}/groups
func (w *RocketmqTopicConsumerGroupsDSWrapper) ListConsumerGroupOfTopic() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dmsv2")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/instances/{instance_id}/topics/{topic}/groups"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	uri = strings.ReplaceAll(uri, "{topic}", w.Get("topic_name").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("groups", "offset", "limit", 0).
		OkCode(200).
		Request().
		Result()
}

func (w *RocketmqTopicConsumerGroupsDSWrapper) listConsumerGroupOfTopicToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("groups", body.Get("groups").Value()),
	)
	return mErr.ErrorOrNil()
}
