package sql

import (
	"testing"

	"github.com/databrickslabs/terraform-provider-databricks/qa"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceSQLGlobalConfigCreateDefault(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "PUT",
				Resource: "/api/2.0/sql/config/endpoints",
				ExpectedRequest: map[string]interface{}{
					"data_access_config":        []interface{}{},
					"enable_serverless_compute": false,
					"security_policy":           "DATA_ACCESS_CONTROL",
				},
			},
			{
				Method:       "GET",
				Resource:     "/api/2.0/sql/config/endpoints",
				ReuseRequest: true,
				Response: GlobalConfigForRead{
					SecurityPolicy: "DATA_ACCESS_CONTROL",
				},
			},
		},
		Resource: ResourceSQLGlobalConfig(),
		Create:   true,
		HCL: `
		`,
	}.Apply(t)
	require.NoError(t, err, err)
	assert.Equal(t, "global", d.Id(), "Id should not be empty")
	assert.Equal(t, "DATA_ACCESS_CONTROL", d.Get("security_policy"))
}

func TestResourceSQLGlobalConfigDelete(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "PUT",
				Resource: "/api/2.0/sql/config/endpoints",
				ExpectedRequest: map[string]interface{}{
					"data_access_config":        []interface{}{},
					"enable_serverless_compute": false,
					"security_policy":           "DATA_ACCESS_CONTROL",
				},
			},
			{
				Method:       "GET",
				Resource:     "/api/2.0/sql/config/endpoints",
				ReuseRequest: true,
				Response: GlobalConfigForRead{
					SecurityPolicy: "DATA_ACCESS_CONTROL",
				},
			},
		},
		Resource: ResourceSQLGlobalConfig(),
		Delete:   true,
		ID:       "global",
		HCL: `
		`,
	}.Apply(t)
	require.NoError(t, err, err)
	assert.Equal(t, "global", d.Id(), "Id should not be empty")
	assert.Equal(t, "DATA_ACCESS_CONTROL", d.Get("security_policy"))
}

func TestResourceSQLGlobalConfigCreateWithData(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "PUT",
				Resource: "/api/2.0/sql/config/endpoints",
				ExpectedRequest: GlobalConfigForRead{
					DataAccessConfig:           []confPair{{Key: "spark.sql.session.timeZone", Value: "UTC"}},
					SqlConfigurationParameters: &repeatedEndpointConfPairs{ConfigPairs: []confPair{{Key: "ANSI_MODE", Value: "true"}}},
					EnableServerlessCompute:    false,
					SecurityPolicy:             "PASSTHROUGH",
					InstanceProfileARN:         "arn:...",
				},
			},
			{
				Method:       "GET",
				Resource:     "/api/2.0/sql/config/endpoints",
				ReuseRequest: true,
				Response: GlobalConfigForRead{
					SecurityPolicy: "PASSTHROUGH",
					DataAccessConfig: []confPair{
						{Key: "spark.sql.session.timeZone", Value: "UTC"},
					},
					InstanceProfileARN: "arn:...",
					SqlConfigurationParameters: &repeatedEndpointConfPairs{
						ConfigPairs: []confPair{
							{Key: "ANSI_MODE", Value: "true"},
						},
					},
				},
			},
		},
		Resource: ResourceSQLGlobalConfig(),
		Create:   true,
		State: map[string]interface{}{
			"security_policy":      "PASSTHROUGH",
			"instance_profile_arn": "arn:...",
			"data_access_config": map[string]interface{}{
				"spark.sql.session.timeZone": "UTC",
			},
			"sql_config_params": map[string]interface{}{
				"ANSI_MODE": "true",
			},
		},
	}.Apply(t)
	require.NoError(t, err, err)
	assert.Equal(t, "global", d.Id(), "Id should not be empty")
	assert.Equal(t, "PASSTHROUGH", d.Get("security_policy"))
}

func TestResourceSQLGlobalConfigCreateError(t *testing.T) {
	_, err := qa.ResourceFixture{
		Resource: ResourceSQLGlobalConfig(),
		Create:   true,
		Azure:    true,
		State: map[string]interface{}{
			"security_policy":      "PASSTHROUGH",
			"instance_profile_arn": "arn:...",
			"data_access_config": map[string]interface{}{
				"spark.sql.session.timeZone": "UTC",
			},
		},
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "can't use instance_profile_arn outside of AWS")
}