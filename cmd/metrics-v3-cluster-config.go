// Copyright (c) 2015-2024 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import "context"

const (
	configWriteQuorum    = "write_quorum"
	configRRSParity      = "rrs_parity"
	configStandardParity = "standard_parity"
)

var (
	configWriteQuorumMD = NewGaugeMD(configWriteQuorum,
		"Maximum write quorum across all pools and sets")
	configRRSParityMD = NewGaugeMD(configRRSParity,
		"Reduced redundancy storage class parity")
	configStandardParityMD = NewGaugeMD(configStandardParity,
		"Standard storage class parity")
)

// loadClusterConfigMetrics - `MetricsLoaderFn` for cluster config
// such as standard and RRS parity.
func loadClusterConfigMetrics(ctx context.Context, m MetricValues, c *metricsCache) error {
	clusterDriveMetrics, _ := c.clusterDriveMetrics.Get()
	m.Set(configStandardParity, float64(clusterDriveMetrics.storageInfo.Backend.StandardSCParity))
	m.Set(configRRSParity, float64(clusterDriveMetrics.storageInfo.Backend.RRSCParity))

	objLayer := newObjectLayerFn()
	if objLayer != nil {
		result := objLayer.Health(ctx, HealthOptions{})
		m.Set(configWriteQuorum, float64(result.WriteQuorum))
	}

	return nil
}
