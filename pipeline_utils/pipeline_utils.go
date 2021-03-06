// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package pipeline_utils

import (
	"errors"
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/common"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/parts"
	"github.com/couchbase/goxdcr/service_def"
	"strconv"
	"strings"
)

var ErrorNoSourceKV = errors.New("Invalid configuration. No source kv node is found.")

func GetSourceVBListPerPipeline(pipeline common.Pipeline) []uint16 {
	ret := []uint16{}
	sourceNozzles := pipeline.Sources()
	for _, sourceNozzle := range sourceNozzles {
		ret = append(ret, sourceNozzle.(*parts.DcpNozzle).GetVBList()...)
	}
	return ret
}

/**
 * Returns a map of: kvServerNode -> vBuckets that it is responsible for
 */
func GetSourceVBMap(cluster_info_svc service_def.ClusterInfoSvc, xdcr_topology_svc service_def.XDCRCompTopologySvc,
	sourceBucketName string, logger *log.CommonLogger) (kv_vb_map map[string][]uint16, number_of_source_nodes int, err error) {
	kv_vb_map = make(map[string][]uint16)

	server_vbmap, err := cluster_info_svc.GetLocalServerVBucketsMap(xdcr_topology_svc, sourceBucketName)
	if err != nil {
		return
	}

	number_of_source_nodes = len(server_vbmap)

	logger.Debugf("server_vbmap=%v\n", server_vbmap)
	nodes, err := xdcr_topology_svc.MyKVNodes()
	if err != nil {
		logger.Errorf("Failed to get my KV nodes, err=%v\n", err)
		return
	}

	if len(nodes) == 0 {
		err = ErrorNoSourceKV
		return
	}

	for _, node := range nodes {
		var kvaddr string = ""
		var vbnos []uint16
		// iterate through serverVBMap and look for server addr that starts with "kvHost:"
		for kvaddr_iter, vbnos_iter := range server_vbmap {
			if kvaddr_iter == node {
				kvaddr = kvaddr_iter
				vbnos = vbnos_iter
				break
			}
		}
		if kvaddr != "" {
			kv_vb_map[kvaddr] = vbnos
		}
	}
	return
}

// checks if target cluster supports SANs in certificates
func HasSANInCertificateSupport(cluster_info_svc service_def.ClusterInfoSvc, targetClusterRef *metadata.RemoteClusterReference) (bool, error) {
	return cluster_info_svc.IsClusterCompatible(targetClusterRef, base.VersionForSANInCertificateSupport)
}

func GetElementIdFromName(pipeline common.Pipeline, name string) string {
	return pipeline.Topic() + "_" + name
}

func GetElementIdFromNameAndIndex(pipeline common.Pipeline, name string, index int) string {
	return pipeline.Topic() + "_" + name + "_" + strconv.Itoa(index)
}

// as the name indicates, this method works only on element ids with index at the end
func GetElementNameFromIdWithIndex(id string) string {
	parts := strings.Split(id, "_")
	if len(parts) > 2 {
		return parts[len(parts)-2]
	}
	return ""
}

func RegisterAsyncComponentEventHandler(listenerMap map[string]common.AsyncComponentEventListener, listenerName string, handler common.AsyncComponentEventHandler) {
	for id, listener := range listenerMap {
		if listenerName == GetElementNameFromIdWithIndex(id) {
			listener.RegisterComponentEventHandler(handler)
		}
	}

}

func IsPipelineUsingCapi(pipeline common.Pipeline) bool {
	isCapi := false
	for _, target := range pipeline.Targets() {
		_, isCapi = target.(*parts.CapiNozzle)
		break
	}
	return isCapi
}

// check if the specified state indicates that the pipeline is still running
// many components, e.g., topology change detector, use this to decide if they need to abort operation
func IsPipelineRunning(state common.PipelineState) bool {
	return state == common.Pipeline_Initial || state == common.Pipeline_Starting || state == common.Pipeline_Running
}

func IsPipelineStopping(state common.PipelineState) bool {
	return state == common.Pipeline_Stopping || state == common.Pipeline_Stopped
}

func DelCheckpointsDocWithRetry(checkpoints_svc service_def.CheckpointsService, replicationId string, vbno uint16, maxRetry int, logger *log.CommonLogger) error {
	logger.Infof("DelCheckpointsDoc for replication %v and vbno %v with maxRetry=%v...", replicationId, vbno, maxRetry)
	numRetry := 0
	for {
		err := checkpoints_svc.DelCheckpointsDoc(replicationId, vbno)
		if err == nil {
			return nil
		}
		if numRetry >= maxRetry {
			errMsg := fmt.Sprintf("Error deleting checkpoint doc for replication %v and vb %v after %v retries", replicationId, vbno, numRetry)
			logger.Warn(errMsg)
			return errors.New(errMsg)
		}
		numRetry++
	}
	return nil
}

// get first seen xattr seqno for each vbucket in pipeline
func GetXattrSeqnos(pipeline common.Pipeline) map[uint16]uint64 {
	ret := make(map[uint16]uint64)
	sourceNozzles := pipeline.Sources()
	for _, sourceNozzle := range sourceNozzles {
		xattr_seqnos := sourceNozzle.(*parts.DcpNozzle).GetXattrSeqnos()
		for vbno, xattr_seqno := range xattr_seqnos {
			ret[vbno] = xattr_seqno
		}
	}
	return ret
}

// get vbnos of checkpoint docs for specified replicationId
func GetVbnosFromCheckpointDocsWithRetry(checkpoints_svc service_def.CheckpointsService, replicationId string, maxRetry int, logger *log.CommonLogger) ([]uint16, error) {
	var vbnos []uint16
	var err error
	numRetry := 0
	for {
		vbnos, err = checkpoints_svc.GetVbnosFromCheckpointDocs(replicationId)
		if err == nil {
			break
		}
		if numRetry >= maxRetry {
			errMsg := fmt.Sprintf("Error getting vbucket numbers for %v after %v retries. err=%v", replicationId, numRetry, err)
			logger.Warn(errMsg)
			return nil, errors.New(errMsg)
		}
		numRetry++
	}
	return vbnos, nil
}

// set replication spec deleted flag in checkpoint docs
func DelCheckpointsDocsWithRetry(checkpoints_svc service_def.CheckpointsService, replicationId string, maxRetry int, logger *log.CommonLogger) error {
	logger.Infof("DelCheckpointsDocs for replication %v with maxRety=%v", replicationId, maxRetry)

	vbnos, err := GetVbnosFromCheckpointDocsWithRetry(checkpoints_svc, replicationId, maxRetry, logger)
	if err != nil {
		return err
	}

	logger.Infof("DelCheckpointsDocs for replication %v with vbs=%v", replicationId, vbnos)

	numRetry := 0
	for {
		failed_vbnos := make([]uint16, 0)
		for _, vbno := range vbnos {
			err = checkpoints_svc.DelCheckpointsDoc(replicationId, vbno)
			if err != nil {
				failed_vbnos = append(failed_vbnos, vbno)
			}
		}
		if len(failed_vbnos) == 0 {
			// no failed vbs, return success
			logger.Infof("DelCheckpointsDocs for replication %v succeeded.", replicationId)
			return nil
		}
		if numRetry >= maxRetry {
			errMsg := fmt.Sprintf("DelCheckpointsDocs for replication %v failed after %v retries. failed vbs=%v", replicationId, numRetry, failed_vbnos)
			logger.Warn(errMsg)
			return errors.New(errMsg)
		}
		numRetry++
		// retry failed vbnos only
		vbnos = failed_vbnos
	}

	return nil
}
