// Code generated by "stringer -type=batchJobMetric -trimprefix=batchJobMetric batch-handlers.go"; DO NOT EDIT.

package cmd

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[batchJobMetricReplication-0]
	_ = x[batchJobMetricKeyRotation-1]
	_ = x[batchJobMetricExpire-2]
}

const _batchJobMetric_name = "ReplicationKeyRotationExpire"

var _batchJobMetric_index = [...]uint8{0, 11, 22, 28}

func (i batchJobMetric) String() string {
	if i >= batchJobMetric(len(_batchJobMetric_index)-1) {
		return "batchJobMetric(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _batchJobMetric_name[_batchJobMetric_index[i]:_batchJobMetric_index[i+1]]
}
