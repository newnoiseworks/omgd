package utils

import (
	"testing"
)

func TestGCPZoneToRegion(t *testing.T) {
	zone := "us-east1-c"

	region := GCPZoneToRegion(zone)

	if region != "us-east1" {
		LogError("GCP Zone not converting to region properly")
		testLogComparison("us-east1", region)
		t.Fail()
	}
}
