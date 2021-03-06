package mocks

import base "github.com/couchbase/goxdcr/base"
import common "github.com/couchbase/goxdcr/common"

import log "github.com/couchbase/goxdcr/log"
import metadata "github.com/couchbase/goxdcr/metadata"
import mock "github.com/stretchr/testify/mock"
import parts "github.com/couchbase/goxdcr/parts"
import time "time"

// XDCRFactoryIface is an autogenerated mock type for the XDCRFactoryIface type
type XDCRFactoryIface struct {
	mock.Mock
}

// constructCAPINozzle provides a mock function with given fields: topic, username, password, certificate, vbList, vbCouchApiBaseMap, nozzle_index, logger_ctx
func (_m *XDCRFactoryIface) constructCAPINozzle(topic string, username string, password string, certificate []byte, vbList []uint16, vbCouchApiBaseMap map[uint16]string, nozzle_index int, logger_ctx *log.LoggerContext) (common.Nozzle, error) {
	ret := _m.Called(topic, username, password, certificate, vbList, vbCouchApiBaseMap, nozzle_index, logger_ctx)

	var r0 common.Nozzle
	if rf, ok := ret.Get(0).(func(string, string, string, []byte, []uint16, map[uint16]string, int, *log.LoggerContext) common.Nozzle); ok {
		r0 = rf(topic, username, password, certificate, vbList, vbCouchApiBaseMap, nozzle_index, logger_ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Nozzle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, []byte, []uint16, map[uint16]string, int, *log.LoggerContext) error); ok {
		r1 = rf(topic, username, password, certificate, vbList, vbCouchApiBaseMap, nozzle_index, logger_ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructOutgoingNozzles provides a mock function with given fields: spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx
func (_m *XDCRFactoryIface) constructOutgoingNozzles(spec *metadata.ReplicationSpecification, kv_vb_map map[string][]uint16, sourceCRMode base.ConflictResolutionMode, targetBucketInfo map[string]interface{}, targetClusterRef *metadata.RemoteClusterReference, logger_ctx *log.LoggerContext) (map[string]common.Nozzle, map[uint16]string, map[string][]uint16, string, string, bool, error) {
	ret := _m.Called(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)

	var r0 map[string]common.Nozzle
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, map[string][]uint16, base.ConflictResolutionMode, map[string]interface{}, *metadata.RemoteClusterReference, *log.LoggerContext) map[string]common.Nozzle); ok {
		r0 = rf(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]common.Nozzle)
		}
	}

	var r1 map[uint16]string
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification, map[string][]uint16, base.ConflictResolutionMode, map[string]interface{}, *metadata.RemoteClusterReference, *log.LoggerContext) map[uint16]string); ok {
		r1 = rf(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[uint16]string)
		}
	}

	var r2 map[string][]uint16
	if rf, ok := ret.Get(2).(func(*metadata.ReplicationSpecification, map[string][]uint16, base.ConflictResolutionMode, map[string]interface{}, *metadata.RemoteClusterReference, *log.LoggerContext) map[string][]uint16); ok {
		r2 = rf(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(map[string][]uint16)
		}
	}

	var r3 string
	if rf, ok := ret.Get(3).(func(*metadata.ReplicationSpecification, map[string][]uint16, base.ConflictResolutionMode, map[string]interface{}, *metadata.RemoteClusterReference, *log.LoggerContext) string); ok {
		r3 = rf(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)
	} else {
		r3 = ret.Get(3).(string)
	}

	var r4 string
	if rf, ok := ret.Get(4).(func(*metadata.ReplicationSpecification, map[string][]uint16, base.ConflictResolutionMode, map[string]interface{}, *metadata.RemoteClusterReference, *log.LoggerContext) string); ok {
		r4 = rf(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)
	} else {
		r4 = ret.Get(4).(string)
	}

	var r5 bool
	if rf, ok := ret.Get(5).(func(*metadata.ReplicationSpecification, map[string][]uint16, base.ConflictResolutionMode, map[string]interface{}, *metadata.RemoteClusterReference, *log.LoggerContext) bool); ok {
		r5 = rf(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)
	} else {
		r5 = ret.Get(5).(bool)
	}

	var r6 error
	if rf, ok := ret.Get(6).(func(*metadata.ReplicationSpecification, map[string][]uint16, base.ConflictResolutionMode, map[string]interface{}, *metadata.RemoteClusterReference, *log.LoggerContext) error); ok {
		r6 = rf(spec, kv_vb_map, sourceCRMode, targetBucketInfo, targetClusterRef, logger_ctx)
	} else {
		r6 = ret.Error(6)
	}

	return r0, r1, r2, r3, r4, r5, r6
}

// constructRouter provides a mock function with given fields: id, spec, downStreamParts, vbNozzleMap, sourceCRMode, logger_ctx
func (_m *XDCRFactoryIface) constructRouter(id string, spec *metadata.ReplicationSpecification, downStreamParts map[string]common.Part, vbNozzleMap map[uint16]string, sourceCRMode base.ConflictResolutionMode, logger_ctx *log.LoggerContext) (*parts.Router, error) {
	ret := _m.Called(id, spec, downStreamParts, vbNozzleMap, sourceCRMode, logger_ctx)

	var r0 *parts.Router
	if rf, ok := ret.Get(0).(func(string, *metadata.ReplicationSpecification, map[string]common.Part, map[uint16]string, base.ConflictResolutionMode, *log.LoggerContext) *parts.Router); ok {
		r0 = rf(id, spec, downStreamParts, vbNozzleMap, sourceCRMode, logger_ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parts.Router)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *metadata.ReplicationSpecification, map[string]common.Part, map[uint16]string, base.ConflictResolutionMode, *log.LoggerContext) error); ok {
		r1 = rf(id, spec, downStreamParts, vbNozzleMap, sourceCRMode, logger_ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructSettingsForCapiNozzle provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructSettingsForCapiNozzle(pipeline common.Pipeline, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, map[string]interface{}) error); ok {
		r1 = rf(pipeline, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructSettingsForCheckpointManager provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructSettingsForCheckpointManager(pipeline common.Pipeline, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, map[string]interface{}) error); ok {
		r1 = rf(pipeline, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructSettingsForDcpNozzle provides a mock function with given fields: pipeline, part, settings
func (_m *XDCRFactoryIface) constructSettingsForDcpNozzle(pipeline common.Pipeline, part *parts.DcpNozzle, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, part, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, *parts.DcpNozzle, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, part, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, *parts.DcpNozzle, map[string]interface{}) error); ok {
		r1 = rf(pipeline, part, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructSettingsForStatsManager provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructSettingsForStatsManager(pipeline common.Pipeline, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, map[string]interface{}) error); ok {
		r1 = rf(pipeline, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructSettingsForSupervisor provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructSettingsForSupervisor(pipeline common.Pipeline, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, map[string]interface{}) error); ok {
		r1 = rf(pipeline, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructSettingsForXmemNozzle provides a mock function with given fields: pipeline, part, targetClusterRef, settings, ssl_port_map, isSSLOverMem
func (_m *XDCRFactoryIface) constructSettingsForXmemNozzle(pipeline common.Pipeline, part common.Part, targetClusterRef *metadata.RemoteClusterReference, settings map[string]interface{}, ssl_port_map map[string]uint16, isSSLOverMem bool) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, part, targetClusterRef, settings, ssl_port_map, isSSLOverMem)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.Part, *metadata.RemoteClusterReference, map[string]interface{}, map[string]uint16, bool) map[string]interface{}); ok {
		r0 = rf(pipeline, part, targetClusterRef, settings, ssl_port_map, isSSLOverMem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, common.Part, *metadata.RemoteClusterReference, map[string]interface{}, map[string]uint16, bool) error); ok {
		r1 = rf(pipeline, part, targetClusterRef, settings, ssl_port_map, isSSLOverMem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructSourceNozzles provides a mock function with given fields: spec, topic, logger_ctx
func (_m *XDCRFactoryIface) constructSourceNozzles(spec *metadata.ReplicationSpecification, topic string, logger_ctx *log.LoggerContext) (map[string]common.Nozzle, map[string][]uint16, error) {
	ret := _m.Called(spec, topic, logger_ctx)

	var r0 map[string]common.Nozzle
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, string, *log.LoggerContext) map[string]common.Nozzle); ok {
		r0 = rf(spec, topic, logger_ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]common.Nozzle)
		}
	}

	var r1 map[string][]uint16
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification, string, *log.LoggerContext) map[string][]uint16); ok {
		r1 = rf(spec, topic, logger_ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string][]uint16)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*metadata.ReplicationSpecification, string, *log.LoggerContext) error); ok {
		r2 = rf(spec, topic, logger_ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// constructUpdateSettingsForCapiNozzle provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructUpdateSettingsForCapiNozzle(pipeline common.Pipeline, settings map[string]interface{}) map[string]interface{} {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// constructUpdateSettingsForCheckpointManager provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructUpdateSettingsForCheckpointManager(pipeline common.Pipeline, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, map[string]interface{}) error); ok {
		r1 = rf(pipeline, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructUpdateSettingsForStatsManager provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructUpdateSettingsForStatsManager(pipeline common.Pipeline, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, map[string]interface{}) error); ok {
		r1 = rf(pipeline, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructUpdateSettingsForSupervisor provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructUpdateSettingsForSupervisor(pipeline common.Pipeline, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, map[string]interface{}) error); ok {
		r1 = rf(pipeline, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// constructUpdateSettingsForXmemNozzle provides a mock function with given fields: pipeline, settings
func (_m *XDCRFactoryIface) constructUpdateSettingsForXmemNozzle(pipeline common.Pipeline, settings map[string]interface{}) map[string]interface{} {
	ret := _m.Called(pipeline, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// constructXMEMNozzle provides a mock function with given fields: topic, kvaddr, sourceBucketName, targetBucketName, username, password, nozzle_index, connPoolSize, sourceCRMode, targetBucketInfo, logger_ctx
func (_m *XDCRFactoryIface) constructXMEMNozzle(topic string, kvaddr string, sourceBucketName string, targetBucketName string, username string, password string, nozzle_index int, connPoolSize int, sourceCRMode base.ConflictResolutionMode, targetBucketInfo map[string]interface{}, logger_ctx *log.LoggerContext) common.Nozzle {
	ret := _m.Called(topic, kvaddr, sourceBucketName, targetBucketName, username, password, nozzle_index, connPoolSize, sourceCRMode, targetBucketInfo, logger_ctx)

	var r0 common.Nozzle
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, string, int, int, base.ConflictResolutionMode, map[string]interface{}, *log.LoggerContext) common.Nozzle); ok {
		r0 = rf(topic, kvaddr, sourceBucketName, targetBucketName, username, password, nozzle_index, connPoolSize, sourceCRMode, targetBucketInfo, logger_ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Nozzle)
		}
	}

	return r0
}

// filterVBList provides a mock function with given fields: targetkvVBList, kv_vb_map
func (_m *XDCRFactoryIface) filterVBList(targetkvVBList []uint16, kv_vb_map map[string][]uint16) []uint16 {
	ret := _m.Called(targetkvVBList, kv_vb_map)

	var r0 []uint16
	if rf, ok := ret.Get(0).(func([]uint16, map[string][]uint16) []uint16); ok {
		r0 = rf(targetkvVBList, kv_vb_map)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint16)
		}
	}

	return r0
}

// getOutNozzleType provides a mock function with given fields: targetClusterRef, spec
func (_m *XDCRFactoryIface) getOutNozzleType(targetClusterRef *metadata.RemoteClusterReference, spec *metadata.ReplicationSpecification) (base.XDCROutgoingNozzleType, error) {
	ret := _m.Called(targetClusterRef, spec)

	var r0 base.XDCROutgoingNozzleType
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference, *metadata.ReplicationSpecification) base.XDCROutgoingNozzleType); ok {
		r0 = rf(targetClusterRef, spec)
	} else {
		r0 = ret.Get(0).(base.XDCROutgoingNozzleType)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*metadata.RemoteClusterReference, *metadata.ReplicationSpecification) error); ok {
		r1 = rf(targetClusterRef, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// getTargetTimeoutEstimate provides a mock function with given fields: topic
func (_m *XDCRFactoryIface) getTargetTimeoutEstimate(topic string) time.Duration {
	ret := _m.Called(topic)

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func(string) time.Duration); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// partId provides a mock function with given fields: prefix, topic, kvaddr, index
func (_m *XDCRFactoryIface) partId(prefix string, topic string, kvaddr string, index int) string {
	ret := _m.Called(prefix, topic, kvaddr, index)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, int) string); ok {
		r0 = rf(prefix, topic, kvaddr, index)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// registerAsyncListenersOnSources provides a mock function with given fields: pipeline, logger_ctx
func (_m *XDCRFactoryIface) registerAsyncListenersOnSources(pipeline common.Pipeline, logger_ctx *log.LoggerContext) {
	_m.Called(pipeline, logger_ctx)
}

// registerAsyncListenersOnTargets provides a mock function with given fields: pipeline, logger_ctx
func (_m *XDCRFactoryIface) registerAsyncListenersOnTargets(pipeline common.Pipeline, logger_ctx *log.LoggerContext) {
	_m.Called(pipeline, logger_ctx)
}

// registerServices provides a mock function with given fields: pipeline, logger_ctx, kv_vb_map, targetUserName, targetPassword, targetBucketName, target_kv_vb_map, targetClusterRef, targetHasRBACSupport
func (_m *XDCRFactoryIface) registerServices(pipeline common.Pipeline, logger_ctx *log.LoggerContext, kv_vb_map map[string][]uint16, targetUserName string, targetPassword string, targetBucketName string, target_kv_vb_map map[string][]uint16, targetClusterRef *metadata.RemoteClusterReference, targetHasRBACSupport bool) error {
	ret := _m.Called(pipeline, logger_ctx, kv_vb_map, targetUserName, targetPassword, targetBucketName, target_kv_vb_map, targetClusterRef, targetHasRBACSupport)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline, *log.LoggerContext, map[string][]uint16, string, string, string, map[string][]uint16, *metadata.RemoteClusterReference, bool) error); ok {
		r0 = rf(pipeline, logger_ctx, kv_vb_map, targetUserName, targetPassword, targetBucketName, target_kv_vb_map, targetClusterRef, targetHasRBACSupport)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckpointBeforeStop provides a mock function with given fields: pipeline
func (_m *XDCRFactoryIface) CheckpointBeforeStop(pipeline common.Pipeline) error {
	ret := _m.Called(pipeline)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(pipeline)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConstructSSLPortMap provides a mock function with given fields: targetClusterRef, spec
func (_m *XDCRFactoryIface) ConstructSSLPortMap(targetClusterRef *metadata.RemoteClusterReference, spec *metadata.ReplicationSpecification) (map[string]uint16, bool, error) {
	ret := _m.Called(targetClusterRef, spec)

	var r0 map[string]uint16
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference, *metadata.ReplicationSpecification) map[string]uint16); ok {
		r0 = rf(targetClusterRef, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]uint16)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(*metadata.RemoteClusterReference, *metadata.ReplicationSpecification) bool); ok {
		r1 = rf(targetClusterRef, spec)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*metadata.RemoteClusterReference, *metadata.ReplicationSpecification) error); ok {
		r2 = rf(targetClusterRef, spec)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ConstructSettingsForPart provides a mock function with given fields: pipeline, part, settings, targetClusterRef, ssl_port_map, isSSLOverMem
func (_m *XDCRFactoryIface) ConstructSettingsForPart(pipeline common.Pipeline, part common.Part, settings map[string]interface{}, targetClusterRef *metadata.RemoteClusterReference, ssl_port_map map[string]uint16, isSSLOverMem bool) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, part, settings, targetClusterRef, ssl_port_map, isSSLOverMem)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.Part, map[string]interface{}, *metadata.RemoteClusterReference, map[string]uint16, bool) map[string]interface{}); ok {
		r0 = rf(pipeline, part, settings, targetClusterRef, ssl_port_map, isSSLOverMem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, common.Part, map[string]interface{}, *metadata.RemoteClusterReference, map[string]uint16, bool) error); ok {
		r1 = rf(pipeline, part, settings, targetClusterRef, ssl_port_map, isSSLOverMem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConstructSettingsForService provides a mock function with given fields: pipeline, service, settings
func (_m *XDCRFactoryIface) ConstructSettingsForService(pipeline common.Pipeline, service common.PipelineService, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, service, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.PipelineService, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, service, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, common.PipelineService, map[string]interface{}) error); ok {
		r1 = rf(pipeline, service, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConstructUpdateSettingsForPart provides a mock function with given fields: pipeline, part, settings
func (_m *XDCRFactoryIface) ConstructUpdateSettingsForPart(pipeline common.Pipeline, part common.Part, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, part, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.Part, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, part, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, common.Part, map[string]interface{}) error); ok {
		r1 = rf(pipeline, part, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConstructUpdateSettingsForService provides a mock function with given fields: pipeline, service, settings
func (_m *XDCRFactoryIface) ConstructUpdateSettingsForService(pipeline common.Pipeline, service common.PipelineService, settings map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(pipeline, service, settings)

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func(common.Pipeline, common.PipelineService, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(pipeline, service, settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Pipeline, common.PipelineService, map[string]interface{}) error); ok {
		r1 = rf(pipeline, service, settings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPipeline provides a mock function with given fields: topic, progress_recorder
func (_m *XDCRFactoryIface) NewPipeline(topic string, progress_recorder common.PipelineProgressRecorder) (common.Pipeline, error) {
	ret := _m.Called(topic, progress_recorder)

	var r0 common.Pipeline
	if rf, ok := ret.Get(0).(func(string, common.PipelineProgressRecorder) common.Pipeline); ok {
		r0 = rf(topic, progress_recorder)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Pipeline)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, common.PipelineProgressRecorder) error); ok {
		r1 = rf(topic, progress_recorder)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetStartSeqno provides a mock function with given fields: pipeline
func (_m *XDCRFactoryIface) SetStartSeqno(pipeline common.Pipeline) error {
	ret := _m.Called(pipeline)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(pipeline)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
