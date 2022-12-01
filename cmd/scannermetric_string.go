// Code generated by "stringer -type=scannerMetric -trimprefix=scannerMetric data-scanner-metric.go"; DO NOT EDIT.

package cmd

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[scannerMetricReadMetadata-0]
	_ = x[scannerMetricCheckMissing-1]
	_ = x[scannerMetricSaveUsage-2]
	_ = x[scannerMetricApplyAll-3]
	_ = x[scannerMetricApplyVersion-4]
	_ = x[scannerMetricTierObjSweep-5]
	_ = x[scannerMetricHealCheck-6]
	_ = x[scannerMetricILM-7]
	_ = x[scannerMetricCheckReplication-8]
	_ = x[scannerMetricYield-9]
	_ = x[scannerMetricCleanAbandoned-10]
	_ = x[scannerMetricApplyNonCurrent-11]
	_ = x[scannerMetricStartTrace-12]
	_ = x[scannerMetricScanObject-13]
	_ = x[scannerMetricLastRealtime-14]
	_ = x[scannerMetricScanFolder-15]
	_ = x[scannerMetricScanCycle-16]
	_ = x[scannerMetricScanBucketDrive-17]
	_ = x[scannerMetricLast-18]
}

const _scannerMetric_name = "ReadMetadataCheckMissingSaveUsageApplyAllApplyVersionTierObjSweepHealCheckILMCheckReplicationYieldCleanAbandonedApplyNonCurrentStartTraceScanObjectLastRealtimeScanFolderScanCycleScanBucketDriveLast"

var _scannerMetric_index = [...]uint8{0, 12, 24, 33, 41, 53, 65, 74, 77, 93, 98, 112, 127, 137, 147, 159, 169, 178, 193, 197}

func (i scannerMetric) String() string {
	if i >= scannerMetric(len(_scannerMetric_index)-1) {
		return "scannerMetric(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _scannerMetric_name[_scannerMetric_index[i]:_scannerMetric_index[i+1]]
}
