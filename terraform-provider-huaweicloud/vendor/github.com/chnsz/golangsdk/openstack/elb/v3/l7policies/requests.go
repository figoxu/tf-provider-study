package l7policies

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToL7PolicyCreateMap() (map[string]interface{}, error)
}

type Action string
type RuleType string
type CompareType string

const (
	ActionRedirectToPool     Action = "REDIRECT_TO_POOL"
	ActionRedirectToListener Action = "REDIRECT_TO_LISTENER"
	ActionReject             Action = "REJECT"

	TypeCookie   RuleType = "COOKIE"
	TypeFileType RuleType = "FILE_TYPE"
	TypeHeader   RuleType = "HEADER"
	TypeHostName RuleType = "HOST_NAME"
	TypePath     RuleType = "PATH"

	CompareTypeContains  CompareType = "CONTAINS"
	CompareTypeEndWith   CompareType = "ENDS_WITH"
	CompareTypeEqual     CompareType = "EQUAL_TO"
	CompareTypeRegex     CompareType = "REGEX"
	CompareTypeStartWith CompareType = "STARTS_WITH"
)

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Name of the L7 policy.
	Name string `json:"name,omitempty"`

	// The ID of the listener.
	ListenerID string `json:"listener_id" required:"true"`

	// The L7 policy action. One of REDIRECT_TO_POOL, REDIRECT_TO_URL, or REJECT.
	Action Action `json:"action" required:"true"`

	// The position of this policy on the listener.
	Position int32 `json:"position,omitempty"`

	// The priority of this policy on the listener.
	Priority int32 `json:"priority,omitempty"`

	// A human-readable description for the resource.
	Description string `json:"description,omitempty"`

	// Requests matching this policy will be redirected to the pool with this ID.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolID string `json:"redirect_pool_id,omitempty"`

	// Requests matching this policy will be redirected to this Listener.
	// Only valid if action is REDIRECT_TO_LISTENER.
	RedirectListenerID string `json:"redirect_listener_id,omitempty"`

	// Requests matching this policy will be redirected to the URL.
	// Only valid if action is REDIRECT_TO_URL.
	RedirectUrlConfig *RedirectUrlConfig `json:"redirect_url_config,omitempty"`

	// Requests matching this policy will be redirected to the configuration of the page.
	// Only valid if action is FIXED_RESPONSE.
	FixedResponseConfig *FixedResponseConfig `json:"fixed_response_config,omitempty"`

	// Requests matching this policy will be redirected to the multiple pools. 5 at most
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolsConfig []*RedirectPoolsConfig `json:"redirect_pools_config,omitempty"`

	// Config session persistence between pools that associated with the policy.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolsStickySessionConfig *RedirectPoolsStickySessionConfig `json:"redirect_pools_sticky_session_config,omitempty"`

	// The config of the redirected pool.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolsExtendConfig *RedirectPoolsExtendConfig `json:"redirect_pools_extend_config,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

type RedirectUrlConfig struct {
	// The protocol for redirection.
	Protocol string `json:"protocol,omitempty"`

	// The host name that requests are redirected to.
	Host string `json:"host,omitempty"`

	// The port that requests are redirected to.
	Port string `json:"port,omitempty"`

	// The path that requests are redirected to.
	Path string `json:"path,omitempty"`

	// The query string set in the URL for redirection.
	Query string `json:"query"`

	// The status code returned after the requests are redirected.
	StatusCode string `json:"status_code" required:"true"`

	// The list of request header parameters to be added
	InsertHeadersConfig *InsertHeadersConfig `json:"insert_headers_config,omitempty"`

	// The list of request header parameters to be removed
	RemoveHeadersConfig *RemoveHeadersConfig `json:"remove_headers_config,omitempty"`
}

type FixedResponseConfig struct {
	// The fixed HTTP status code configured in the forwarding rule.
	StatusCode string `json:"status_code" required:"true"`

	// The format of the response body.
	ContentType string `json:"content_type,omitempty"`

	// The content of the response message body.
	MessageBody string `json:"message_body"`

	// The list of request header parameters to be added
	InsertHeadersConfig *InsertHeadersConfig `json:"insert_headers_config,omitempty"`

	// The list of request header parameters to be removed
	RemoveHeadersConfig *RemoveHeadersConfig `json:"remove_headers_config,omitempty"`

	// The limit config of the policy
	TrafficLimitConfig *TrafficLimitConfig `json:"traffic_limit_config,omitempty"`
}

type RedirectPoolsConfig struct {
	// The pool ID
	PoolId string `json:"pool_id" required:"true"`

	// The weight of the pool
	Weight int `json:"weight"`
}

type RedirectPoolsStickySessionConfig struct {
	// Whether enable config session persistence between pools
	Enable bool `json:"enable"`

	// The timeout of the session persistence
	Timeout int `json:"timeout"`
}

type RedirectPoolsExtendConfig struct {
	// Whether the rewriter url enable
	RewriteUrlEnable bool `json:"rewrite_url_enable,omitempty"`

	// The rewriter url config
	RewriteUrlConfig *RewriteUrlConfig `json:"rewrite_url_config,omitempty"`

	// The header parameters to be added
	InsertHeadersConfig *InsertHeadersConfig `json:"insert_headers_config,omitempty"`

	// The header parameters to be removed
	RemoveHeadersConfig *RemoveHeadersConfig `json:"remove_headers_config,omitempty"`

	// The traffic limit config of the policy
	TrafficLimitConfig *TrafficLimitConfig `json:"traffic_limit_config,omitempty"`
}

type RewriteUrlConfig struct {
	// The host of the rewriter url
	Host string `json:"host,omitempty"`

	// The path that requests are redirected to.
	Path string `json:"path,omitempty"`

	// The query string set in the URL for redirection.
	Query string `json:"query"`
}

type InsertHeadersConfig struct {
	// The list of request header parameters to be added
	Configs []*InsertHeaderConfig `json:"configs" required:"true"`
}

type InsertHeaderConfig struct {
	// The parameter name of the added request header
	Key string `json:"key" required:"true"`

	// The value type of the parameter
	ValueType string `json:"value_type" required:"true"`

	// The value of the parameter
	Value string `json:"value" required:"true"`
}

type RemoveHeadersConfig struct {
	// The list of request header parameters to be removed
	Configs []*RemoveHeaderConfig `json:"configs" required:"true"`
}

type RemoveHeaderConfig struct {
	// The parameter name of the removed request header
	Key string `json:"key" required:"true"`
}

type TrafficLimitConfig struct {
	// The overall qps of the policy
	Qps int `json:"qps"`

	// The single source qps of the policy
	PerSourceIpQps int `json:"per_source_ip_qps"`

	// The qps buffer
	Burst int `json:"burst"`
}

// ToL7PolicyCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToL7PolicyCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "l7policy")
}

// Create accepts a CreateOpts struct and uses the values to create a new l7policy.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToL7PolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToL7PolicyListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API.
type ListOpts struct {
	Name               string `q:"name"`
	Description        string `q:"description"`
	ListenerID         string `q:"listener_id"`
	Action             string `q:"action"`
	TenantID           string `q:"tenant_id"`
	RedirectPoolID     string `q:"redirect_pool_id"`
	RedirectListenerID string `q:"redirect_listener_id"`
	Position           int32  `q:"position"`
	AdminStateUp       bool   `q:"admin_state_up"`
	ID                 string `q:"id"`
	Limit              int    `q:"limit"`
	Marker             string `q:"marker"`
	SortKey            string `q:"sort_key"`
	SortDir            string `q:"sort_dir"`
}

// ToL7PolicyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToL7PolicyListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// l7policies. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those l7policies that are owned by the
// project who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToL7PolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return L7PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a particular l7policy based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Delete will permanently delete a particular l7policy based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToL7PolicyUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Name of the L7 policy, empty string is allowed.
	Name *string `json:"name,omitempty"`

	// The L7 policy action. One of REDIRECT_TO_POOL, REDIRECT_TO_URL, or REJECT.
	Action Action `json:"action,omitempty"`

	// The position of this policy on the listener.
	Position int32 `json:"position,omitempty"`

	// The priority of this policy on the listener.
	Priority int32 `json:"priority,omitempty"`

	// A human-readable description for the resource, empty string is allowed.
	Description *string `json:"description,omitempty"`

	// Requests matching this policy will be redirected to the pool with this ID.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolID *string `json:"redirect_pool_id,omitempty"`

	// Requests matching this policy will be redirected to this LISTENER.
	// Only valid if action is REDIRECT_TO_LISTENER.
	RedirectListenerID *string `json:"redirect_listener_id,omitempty"`

	// Requests matching this policy will be redirected to the URL.
	// Only valid if action is REDIRECT_TO_URL.
	RedirectUrlConfig *RedirectUrlConfig `json:"redirect_url_config,omitempty"`

	// Requests matching this policy will be redirected to the configuration of the page.
	// Only valid if action is FIXED_RESPONSE.
	FixedResponseConfig *FixedResponseConfig `json:"fixed_response_config,omitempty"`

	// Requests matching this policy will be redirected to the multiple pools. 5 at most
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolsConfig []*RedirectPoolsConfig `json:"redirect_pools_config,omitempty"`

	// Config session persistence between pools that associated with the policy.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolsStickySessionConfig *RedirectPoolsStickySessionConfig `json:"redirect_pools_sticky_session_config,omitempty"`

	// The config of the redirected pool.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolsExtendConfig *RedirectPoolsExtendConfig `json:"redirect_pools_extend_config,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

// ToL7PolicyUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToL7PolicyUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "l7policy")
	if err != nil {
		return nil, err
	}

	m := b["l7policy"].(map[string]interface{})

	if m["redirect_pool_id"] == "" {
		m["redirect_pool_id"] = nil
	}

	if m["redirect_url"] == "" {
		m["redirect_url"] = nil
	}

	return b, nil
}

// Update allows l7policy to be updated.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToL7PolicyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// CreateRuleOpts is the common options struct used in this package's CreateRule
// operation.
type CreateRuleOpts struct {
	// The L7 rule type. One of HOST_NAME, PATH, METHOD, HEADER, QUERY_STRING, or SOURCE_IP.
	RuleType RuleType `json:"type" required:"true"`

	// The comparison type for the L7 rule. One of EQUAL_TO, REGEX, or STARTS_WITH.
	CompareType CompareType `json:"compare_type" required:"true"`

	// The value to use for the comparison.
	Value string `json:"value,omitempty"`

	// TenantID is the UUID of the tenant who owns the rule in octavia.
	// Only administrative users can specify a project UUID other than their own.
	TenantID string `json:"tenant_id,omitempty"`

	// The key to use for the comparison. For example, the name of the cookie to evaluate.
	Key string `json:"key,omitempty"`

	// When true the logic of the rule is inverted. For example, with invert true,
	// equal to would become not equal to. Default is false.
	Invert bool `json:"invert,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// The matching conditions of the forwarding rule.
	// This parameter is available only when enhance_l7policy_enable of the listener is set to true.
	Conditions []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	// The key of the match item.
	Key string `json:"key,omitempty"`

	// The value of the match item.
	Value string `json:"value" required:"true"`
}

// ToRuleCreateMap builds a request body from CreateRuleOpts.
func (opts CreateRuleOpts) ToRuleCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "rule")
}

// CreateRule will create and associate a Rule with a particular L7Policy.
func CreateRule(c *golangsdk.ServiceClient, policyID string, opts CreateRuleOpts) (r CreateRuleResult) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(ruleRootURL(c, policyID), b, &r.Body, nil)
	return
}

// ListRulesOptsBuilder allows extensions to add additional parameters to the
// ListRules request.
type ListRulesOptsBuilder interface {
	ToRulesListQuery() (string, error)
}

// ListRulesOpts allows the filtering and sorting of paginated collections
// through the API.
type ListRulesOpts struct {
	RuleType     RuleType    `q:"type"`
	TenantID     string      `q:"tenant_id"`
	CompareType  CompareType `q:"compare_type"`
	Value        string      `q:"value"`
	Key          string      `q:"key"`
	Invert       bool        `q:"invert"`
	AdminStateUp bool        `q:"admin_state_up"`
	ID           string      `q:"id"`
	Limit        int         `q:"limit"`
	Marker       string      `q:"marker"`
	SortKey      string      `q:"sort_key"`
	SortDir      string      `q:"sort_dir"`
}

// ToRulesListQuery formats a ListOpts into a query string.
func (opts ListRulesOpts) ToRulesListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// ListRules returns a Pager which allows you to iterate over a collection of
// rules. It accepts a ListRulesOptsBuilder, which allows you to filter and
// sort the returned collection for greater efficiency.
//
// Default policy settings return only those rules that are owned by the
// project who submits the request, unless an admin user submits the request.
func ListRules(c *golangsdk.ServiceClient, policyID string, opts ListRulesOptsBuilder) pagination.Pager {
	url := ruleRootURL(c, policyID)
	if opts != nil {
		query, err := opts.ToRulesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetRule retrieves a particular L7Policy Rule based on its unique ID.
func GetRule(c *golangsdk.ServiceClient, policyID string, ruleID string) (r GetRuleResult) {
	_, r.Err = c.Get(ruleResourceURL(c, policyID, ruleID), &r.Body, nil)
	return
}

// DeleteRule will remove a Rule from a particular L7Policy.
func DeleteRule(c *golangsdk.ServiceClient, policyID string, ruleID string) (r DeleteRuleResult) {
	_, r.Err = c.Delete(ruleResourceURL(c, policyID, ruleID), nil)
	return
}

// UpdateRuleOptsBuilder allows to add additional parameters to the PUT request.
type UpdateRuleOptsBuilder interface {
	ToRuleUpdateMap() (map[string]interface{}, error)
}

// UpdateRuleOpts is the common options struct used in this package's Update
// operation.
type UpdateRuleOpts struct {
	// The L7 rule type. One of HOST_NAME, PATH, METHOD, HEADER, QUERY_STRING, or SOURCE_IP.
	RuleType RuleType `json:"type,omitempty"`

	// The comparison type for the L7 rule. One of EQUAL_TO, REGEX, or STARTS_WITH.
	CompareType CompareType `json:"compare_type,omitempty"`

	// The value to use for the comparison.
	Value string `json:"value,omitempty"`

	// The key to use for the comparison. For example, the name of the cookie to evaluate.
	Key *string `json:"key,omitempty"`

	// When true the logic of the rule is inverted. For example, with invert true,
	// equal to would become not equal to. Default is false.
	Invert *bool `json:"invert,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// The matching conditions of the forwarding rule.
	// This parameter is available only when enhance_l7policy_enable of the listener is set to true.
	Conditions []Condition `json:"conditions,omitempty"`
}

// ToRuleUpdateMap builds a request body from UpdateRuleOpts.
func (opts UpdateRuleOpts) ToRuleUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "rule")
	if err != nil {
		return nil, err
	}

	if m := b["rule"].(map[string]interface{}); m["key"] == "" {
		m["key"] = nil
	}

	return b, nil
}

// UpdateRule allows Rule to be updated.
func UpdateRule(c *golangsdk.ServiceClient, policyID string, ruleID string, opts UpdateRuleOptsBuilder) (r UpdateRuleResult) {
	b, err := opts.ToRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(ruleResourceURL(c, policyID, ruleID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	return
}
