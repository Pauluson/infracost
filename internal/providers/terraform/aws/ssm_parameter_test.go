package aws_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/infracost/infracost/internal/schema"
	"github.com/infracost/infracost/internal/testutil"

	"github.com/infracost/infracost/internal/providers/terraform/tftest"
)

func TestAwsSSMParameterFunction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	tf := `
  resource "aws_ssm_parameter" "advanced" {
    name = "my-advanced-ssm-parameter"
		type = "String"
		value = "Advanced Parameter"
		tier = "Advanced"
  }
	`

	resourceChecks := []testutil.ResourceCheck{
		{
			Name: "aws_ssm_parameter.advanced",
			CostComponentChecks: []testutil.CostComponentCheck{
				{
					Name:             "Parameter storage (advanced)",
					PriceHash:        "d5db437b8b7a6df9c701534aefab452b-1065e83bbc0d4959dda26a1848f3e3eb",
					MonthlyCostCheck: testutil.HourlyPriceMultiplierCheck(decimal.NewFromInt(1)),
				},
				{
					Name:             "API interactions (advanced)",
					PriceHash:        "8857de3489efa197e0f05fdbc54c760f-7c35c68819b19a7ff1d898cc5a198a7f",
					MonthlyCostCheck: testutil.NilMonthlyCostCheck(),
				},
			},
		},
	}

	tftest.ResourceTests(t, tf, schema.NewEmptyUsageMap(), resourceChecks)
}

func TestAwsSSMParameter_usage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	tf := `
  resource "aws_ssm_parameter" "advanced" {
    name = "my-advanced-ssm-parameter"
		type = "String"
		value = "Advanced Parameter"
		tier = "Advanced"
  }
	`

	usage := schema.NewUsageMap(map[string]interface{}{
		"aws_ssm_parameter.advanced": map[string]interface{}{
			"api_throughput_limit":     "advanced",
			"monthly_api_interactions": 100000,
			"parameter_storage_hrs":    600,
		},
	})

	resourceChecks := []testutil.ResourceCheck{
		{
			Name: "aws_ssm_parameter.advanced",
			CostComponentChecks: []testutil.CostComponentCheck{
				{
					Name:             "Parameter storage (advanced)",
					PriceHash:        "d5db437b8b7a6df9c701534aefab452b-1065e83bbc0d4959dda26a1848f3e3eb",
					MonthlyCostCheck: testutil.MonthlyPriceMultiplierCheck(decimal.NewFromInt(600)),
				},
				{
					Name:             "API interactions (advanced)",
					PriceHash:        "8857de3489efa197e0f05fdbc54c760f-7c35c68819b19a7ff1d898cc5a198a7f",
					MonthlyCostCheck: testutil.MonthlyPriceMultiplierCheck(decimal.NewFromInt(100000)),
				},
			},
		},
	}

	tftest.ResourceTests(t, tf, usage, resourceChecks)
}
