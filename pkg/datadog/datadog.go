package datadog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/config"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/helpers"
	"github.com/sirupsen/logrus"
)

func generateBody(item string, timestamp int64) datadogV2.MetricPayload {
	body := datadogV2.MetricPayload{
		Series: []datadogV2.MetricSeries{
			{
				Metric: "dritctl.coverage",
				Type:   datadogV2.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(timestamp),
						Value:     datadog.PtrFloat64(config.DritfctlRunMap[item].Summary.Coverage),
					},
				},
				Tags: []string{
					fmt.Sprintf("runId:%s", item),
					fmt.Sprintf("provider:%s", config.DritfctlRunMap[item].Provider),
					fmt.Sprintf("awsAccount:%s", config.DritfctlRunMap[item].AWSConfig.AWSAccountId),
					fmt.Sprintf("group:%s", config.DritfctlRunMap[item].Group),
				},
			},
			{
				Metric: "dritctl.total_resources",
				Type:   datadogV2.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(timestamp),
						Value:     datadog.PtrFloat64(config.DritfctlRunMap[item].Summary.TotalResources),
					},
				},
				Tags: []string{
					fmt.Sprintf("runId:%s", item),
					fmt.Sprintf("provider:%s", config.DritfctlRunMap[item].Provider),
					fmt.Sprintf("awsAccount:%s", config.DritfctlRunMap[item].AWSConfig.AWSAccountId),
					fmt.Sprintf("group:%s", config.DritfctlRunMap[item].Group),
				},
			},
			{
				Metric: "dritctl.total_changed",
				Type:   datadogV2.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(timestamp),
						Value:     datadog.PtrFloat64(config.DritfctlRunMap[item].Summary.TotalChanged),
					},
				},
				Tags: []string{
					fmt.Sprintf("runId:%s", item),
					fmt.Sprintf("provider:%s", config.DritfctlRunMap[item].Provider),
					fmt.Sprintf("awsAccount:%s", config.DritfctlRunMap[item].AWSConfig.AWSAccountId),
					fmt.Sprintf("group:%s", config.DritfctlRunMap[item].Group),
				},
			},
			{
				Metric: "dritctl.total_unmanaged",
				Type:   datadogV2.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(timestamp),
						Value:     datadog.PtrFloat64(config.DritfctlRunMap[item].Summary.TotalUnmanaged),
					},
				},
				Tags: []string{
					fmt.Sprintf("runId:%s", item),
					fmt.Sprintf("provider:%s", config.DritfctlRunMap[item].Provider),
					fmt.Sprintf("awsAccount:%s", config.DritfctlRunMap[item].AWSConfig.AWSAccountId),
					fmt.Sprintf("group:%s", config.DritfctlRunMap[item].Group),
				},
			},
			{
				Metric: "dritctl.total_managed",
				Type:   datadogV2.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(timestamp),
						Value:     datadog.PtrFloat64(config.DritfctlRunMap[item].Summary.TotalManaged),
					},
				},
				Tags: []string{
					fmt.Sprintf("runId:%s", item),
					fmt.Sprintf("provider:%s", config.DritfctlRunMap[item].Provider),
					fmt.Sprintf("awsAccount:%s", config.DritfctlRunMap[item].AWSConfig.AWSAccountId),
					fmt.Sprintf("group:%s", config.DritfctlRunMap[item].Group),
				},
			},
			{
				Metric: "dritctl.total_deleted",
				Type:   datadogV2.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(timestamp),
						Value:     datadog.PtrFloat64(config.DritfctlRunMap[item].Summary.TotalDeleted),
					},
				},
				Tags: []string{
					fmt.Sprintf("runId:%s", item),
					fmt.Sprintf("provider:%s", config.DritfctlRunMap[item].Provider),
					fmt.Sprintf("awsAccount:%s", config.DritfctlRunMap[item].AWSConfig.AWSAccountId),
					fmt.Sprintf("group:%s", config.DritfctlRunMap[item].Group),
				},
			},
		},
	}
	return body
}

func SendMetrics(item string) {
	if !helpers.EnvExist("DD_SITE") || !helpers.EnvExist("DD_API_KEY") {
		logrus.Error("DD_SITE and DD_API_KEY datadog environment variables are not defined.")
		return
	}
	body := generateBody(item, time.Now().Unix())
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewMetricsApi(apiClient)
	resp, r, err := api.SubmitMetrics(ctx, body, *datadogV2.NewSubmitMetricsOptionalParameters())

	if err != nil {
		logrus.Errorf("Error when calling datadog MetricsApi.SubmitMetrics: %v", err)
		logrus.Errorf("Full HTTP datadog response: %v", r)
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	logrus.Infof("Response from datadog MetricsApi.SubmitMetrics: %s", responseContent)
}
