// Generated by PMS #249
package dws

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

func DataSourceDwsLogicalClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDwsLogicalClustersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specified the cluster ID of the DWS.`,
			},
			"add_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the logical cluster can be added.`,
			},
			"logical_clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All logical clusters that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the logical cluster.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the logical cluster.`,
						},
						"first_logical_cluster": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether it is the first logical cluster.`,
						},
						"cluster_rings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of logical cluster rings.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ring_hosts": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of the cluster hosts.`,
										Elem:        logCluCluRinRinHostsElem(),
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the logical cluster.`,
						},
						"edit_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the logical cluster is allowed to be edited.`,
						},
						"restart_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the logical cluster is allowed to be restarted.`,
						},
						"delete_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the logical cluster is allowed to be deleted.`,
						},
					},
				},
			},
		},
	}
}

// logCluCluRinRinHostsElem
// The Elem of "logical_clusters.cluster_rings.ring_hosts"
func logCluCluRinRinHostsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The host name.`,
			},
			"back_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backend IP address.`,
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of CPU cores.`,
			},
			"memory": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The host memory, in GB.`,
			},
			"disk_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The host disk size, in GB.`,
			},
		},
	}
}

type LogicalClustersDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newLogicalClustersDSWrapper(d *schema.ResourceData, meta interface{}) *LogicalClustersDSWrapper {
	return &LogicalClustersDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDwsLogicalClustersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newLogicalClustersDSWrapper(d, meta)
	lisLogCluRst, err := wrapper.ListLogicalClusters()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listLogicalClustersToSchema(lisLogCluRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/logical-clusters
func (w *LogicalClustersDSWrapper) ListLogicalClusters() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dws")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/clusters/{cluster_id}/logical-clusters"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("logical_clusters", "offset", "limit", 0).
		Request().
		Result()
}

func (w *LogicalClustersDSWrapper) listLogicalClustersToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("add_enable", body.Get("add_enable").Value()),
		d.Set("logical_clusters", schemas.SliceToList(body.Get("logical_clusters"),
			func(logicalClusters gjson.Result) any {
				return map[string]any{
					"id":                    logicalClusters.Get("logical_cluster_id").Value(),
					"name":                  logicalClusters.Get("logical_cluster_name").Value(),
					"first_logical_cluster": logicalClusters.Get("first_logical_cluster").Value(),
					"cluster_rings": schemas.SliceToList(logicalClusters.Get("cluster_rings"),
						func(clusterRings gjson.Result) any {
							return map[string]any{
								"ring_hosts": w.setLogCluCluRinRinHosts(clusterRings),
							}
						},
					),
					"status":         logicalClusters.Get("status").Value(),
					"edit_enable":    logicalClusters.Get("edit_enable").Value(),
					"restart_enable": logicalClusters.Get("restart_enable").Value(),
					"delete_enable":  logicalClusters.Get("delete_enable").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*LogicalClustersDSWrapper) setLogCluCluRinRinHosts(clusterRings gjson.Result) any {
	return schemas.SliceToList(clusterRings.Get("ring_hosts"), func(ringHosts gjson.Result) any {
		return map[string]any{
			"host_name": ringHosts.Get("host_name").Value(),
			"back_ip":   ringHosts.Get("back_ip").Value(),
			"cpu_cores": ringHosts.Get("cpu_cores").Value(),
			"memory":    ringHosts.Get("memory").Value(),
			"disk_size": ringHosts.Get("disk_size").Value(),
		}
	})
}
