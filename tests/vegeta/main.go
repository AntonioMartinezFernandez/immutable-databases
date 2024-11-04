package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	vegeta "github.com/tsenart/vegeta/lib"
)

// Config holds the configuration for the load test.
type Config struct {
	Duration time.Duration
	RateFreq int
	RatePer  time.Duration
}

// TargetConfig represents the configuration for a single target.
type TargetConfig struct {
	URL          string
	Method       string
	Header       map[string][]string
	BodyTemplate string
	ReportPath   string
}

// newDynamicTargeter creates a vegeta.Targeter for a single target configuration.
func newDynamicTargeter(tc TargetConfig) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		// Generate a new UUID for dynamic URL parameter.
		uuid := uuid.New()
		body := fmt.Sprintf(tc.BodyTemplate, uuid.String(), uuid.String())

		tgt.Method = tc.Method
		tgt.URL = tc.URL
		tgt.Header = tc.Header
		tgt.Body = []byte(body)

		return nil
	}
}

// executeAttack runs the attack using the provided configuration, targeter, and report path.
func executeAttack(config Config, targeter vegeta.Targeter, reportPath string) error {
	attacker := vegeta.NewAttacker()
	rate := vegeta.Rate{Freq: config.RateFreq, Per: config.RatePer}
	metrics := &vegeta.Metrics{}

	// Start the attack and collect metrics.
	for res := range attacker.Attack(targeter, rate, config.Duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	// Generate and save the report.
	reporter := vegeta.NewHDRHistogramPlotReporter(metrics)
	f, err := os.Create(reportPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer f.Close()

	if err = reporter.Report(f); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	log.Printf("Report saved to %s", f.Name())
	return nil
}

func main() {
	// Configure the list of target URLs and their settings.
	targetConfigs := []TargetConfig{
		{
			URL:    "http://localhost:3000/api/tracking/events",
			Method: "POST",
			Header: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {"Basic bWFub2xpOmhvbGk="}, // user:manoli password:holi
				"x-api-key":     {"000"},
			},
			BodyTemplate: `{
				"version": 2,
				"operationId": "%s",
				"tenantId": "bab1fc99-df84-4239-9998-957039d515b4",
				"sessionId": "%s",
				"source": "sdk.mobile",
				"family": "ONBOARDING",
				"events": [
					{
						"eventId": "79204f02-0ad8-4665-83a0-62024b2f2b50",
						"clientTimestamp": "2011-10-05T14:48:00.000Z",
						"executionTime": null,
						"payload": {
							"type": "STEP_CHANGE",
							"stepType": "START",
							"stepId": "79204f02-0ad8-4665-83a0-62024b2f3f23",
							"component": {
								"id": "tracking_component",
								"version": "2.0.2-SNAPSHOT"
							}
						}
					}
				]
			}`,
			ReportPath: "./reports/vegeta-report-create-asset-immudb.hgrm",
		},
		// Additional targets can be added here.
	}

	// Set up the load test configuration.
	config := Config{
		Duration: 30 * time.Second,
		RateFreq: 500,
		RatePer:  time.Second,
	}

	// For each target configuration, execute the attack.
	for _, tc := range targetConfigs {
		// Initialize the targeter with the target configuration.
		targeter := newDynamicTargeter(tc)

		// Execute the load test attack for this target.
		if err := executeAttack(config, targeter, tc.ReportPath); err != nil {
			log.Fatalf("Error executing attack for target %s: %v", tc.URL, err)
		}
	}
}
