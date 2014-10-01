// replication manager.

package replication_manager

import (
	"errors"
	"fmt"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr/log"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr/pipeline_manager"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/base"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/factory"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/metadata"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/metadata_svc"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/pipeline_svc"
	"sync"
)

var logger_rm *log.CommonLogger = log.NewLogger("ReplicationManager", log.DefaultLoggerContext)

/************************************
/* struct ReplicationManager
*************************************/
type replicationManager struct {
	metadata_svc             metadata_svc.MetadataSvc
	cluster_info_svc         metadata_svc.ClusterInfoSvc
	xdcr_topology_svc        metadata_svc.XDCRCompTopologySvc
	replication_settings_svc metadata_svc.ReplicationSettingsSvc
	once                     sync.Once
}

var replication_mgr replicationManager

func Initialize(metadata_svc metadata_svc.MetadataSvc,
	cluster_info_svc metadata_svc.ClusterInfoSvc,
	xdcr_topology_svc metadata_svc.XDCRCompTopologySvc,
	replication_settings_svc metadata_svc.ReplicationSettingsSvc) {
	replication_mgr.once.Do(func() {
		replication_mgr.init(metadata_svc, cluster_info_svc, xdcr_topology_svc, replication_settings_svc)
	})
}

func (rm *replicationManager) init(metadataSvc metadata_svc.MetadataSvc,
	clusterSvc metadata_svc.ClusterInfoSvc,
	topologySvc metadata_svc.XDCRCompTopologySvc,
	replicationSettingsSvc metadata_svc.ReplicationSettingsSvc) {
	rm.metadata_svc = metadataSvc
	rm.cluster_info_svc = clusterSvc
	rm.xdcr_topology_svc = topologySvc
	rm.replication_settings_svc = replicationSettingsSvc
	fac := factory.NewXDCRFactory(metadataSvc, clusterSvc, topologySvc, log.DefaultLoggerContext, log.DefaultLoggerContext, rm)
	pipeline_manager.PipelineManager(fac, log.DefaultLoggerContext)

	logger_rm.Info("Replication manager is initialized")

}

func MetadataService() metadata_svc.MetadataSvc {
	return replication_mgr.metadata_svc
}

func ClusterInfoService() metadata_svc.ClusterInfoSvc {
	return replication_mgr.cluster_info_svc
}

func XDCRCompTopologyService() metadata_svc.XDCRCompTopologySvc {
	return replication_mgr.xdcr_topology_svc
}

func ReplicationSettingsService() metadata_svc.ReplicationSettingsSvc {
	return replication_mgr.replication_settings_svc
}

func CreateReplication(sourceClusterUUID string, sourceBucket string, targetClusterUUID, targetBucket string, filterName string, settings map[string]interface{}) (string, error) {
	logger_rm.Infof("Creating replication - sourceCluterUUID=%s, sourceBucket=%s, targetClusterUUID=%s, targetBucket=%s, filterName=%s, settings=%v\n", sourceClusterUUID,
		sourceBucket, targetClusterUUID, targetBucket, filterName, settings)
	spec, err := replication_mgr.createAndPersistReplicationSpec(sourceClusterUUID, sourceBucket, targetClusterUUID, targetBucket, filterName, settings)
	logger_rm.Debugf("replication specification %s is created and persisted\n", spec.Id())
	if err != nil {
		logger_rm.Errorf("%v\n", err)
		return "", err
	}
	_, err = pipeline_manager.StartPipeline(spec.Id(), settings)
	if err == nil {
		logger_rm.Debugf("Pipeline %s is started\n", spec.Id())
		return spec.Id(), err
	} else {
		logger_rm.Errorf("%v\n", err)
		return "", err
	}
}

func PauseReplication(topic string) error {
	logger_rm.Infof("Pausing replication %s\n", topic)
	err := pipeline_manager.StopPipeline(topic)
	if err != nil {
		logger_rm.Errorf("%v\n", err)
		return err
	}
	logger_rm.Debugf("Pipeline %s is stopped", topic)
	//update the replication specification's "active" setting to be false
	spec, err := replication_mgr.metadata_svc.ReplicationSpec(topic)
	if err != nil {
		logger_rm.Errorf("%v\n", err)
		return err
	}

	settings := spec.Settings()
	settings.SetActive(false)
	err = replication_mgr.metadata_svc.SetReplicationSpec(*spec)
	logger_rm.Debugf("Replication specification %s is set to active=false\n", topic)
	return err
}

func ResumeReplication(topic string) error {
	logger_rm.Infof("Resuming replication %s\n", topic)
	spec, err := replication_mgr.metadata_svc.ReplicationSpec(topic)
	if err != nil {
		logger_rm.Errorf("%v\n", err)
		return err
	}

	settings := spec.Settings()
	settingsMap := settings.ToMap()
	_, err = pipeline_manager.StartPipeline(topic, settingsMap)
	if err != nil {
		logger_rm.Errorf("%v\n", err)
		return err
	}
	logger_rm.Debugf("Pipeline %s is started", topic)

	settings.SetActive(true)
	err = replication_mgr.metadata_svc.SetReplicationSpec(*spec)
	if err == nil {
		logger_rm.Debugf("Replication specification %s is updated with active=true\n", topic)
	} else {
		logger_rm.Errorf("%v\n", err)
	}
	return err
}

func DeleteReplication(topic string) error {
	logger_rm.Infof("Deleting replication %s\n", topic)
	err := pipeline_manager.StopPipeline(topic)
	if err != nil {
		logger_rm.Errorf("%v\n", err)
		return err
	}
	logger_rm.Debugf("Pipeline %s is stopped", topic)

	err = replication_mgr.metadata_svc.DelReplicationSpec(topic)
	if err == nil {
		logger_rm.Debugf("Replication specification %s is deleted\n", topic)
	} else {
		logger_rm.Errorf("%v\n", err)
	}
	return err
}

func HandleChangesToReplicationSettings(topic string, settings map[string]interface{}) error {
	// read replication spec with the specified replication id
	replSpec, err := MetadataService().ReplicationSpec(topic)
	if err != nil {
		return err
	}

	// update replication spec with input settings
	replSpec.Settings().UpdateSettingsFromMap(settings)
	err = MetadataService().SetReplicationSpec(*replSpec)

	// TODO implement additional logic, e.g.,
	// 1. reconstruct pipeline when source/targetNozzlePerNode is changed
	// 2. pause pipeline when active is changed from true to false
	// 3. restart pipeline when criteral settings are changed
	return err
}

// get statistics for all running replications
func GetStatistics() (map[string]interface{}, error) {
	// TODO implement
	return nil, nil
}

func (rm *replicationManager) createAndPersistReplicationSpec(sourceClusterUUID, sourceBucket, targetClusterUUID, targetBucket, filterName string, settings map[string]interface{}) (*metadata.ReplicationSpecification, error) {
	spec := metadata.NewReplicationSpecification(sourceClusterUUID, sourceBucket, targetClusterUUID, targetBucket, filterName)
	s, err := metadata.SettingsFromMap(settings)
	if err == nil {
		spec.SetSettings(s)
		return spec, nil
	} else {
		return nil, err
	}
}

func (rm *replicationManager) OnError(topic string, partsError map[string]error) {
	logger_rm.Infof("Pipeline %v reported failure. The following parts are broken: %v\n", topic, partsError)

	//fix the pipeline
	err := fixPipeline(topic)

	if err != nil {
		logger_rm.Infof("The effort of fixing pipeline %v is failed, err=%v; The retry will happen latter\n",
			topic, err)
		//TODO: propagate the error to ns_server

		//TODO: schedule the retry
	}

}

func fixPipeline(topic string) error {
	//pause the replication
	err := PauseReplication(topic)
	if err == nil {
		err = ResumeReplication(topic)
	}
	return err
}

func SetPipelineLogLevel(topic string, levelStr string) error {
	pipeline := pipeline_manager.Pipeline(topic)

	//update the setting
	spec, err := MetadataService().ReplicationSpec(topic)
	if err != nil && spec == nil {
		return errors.New(fmt.Sprintf("Failed to lookup replication specification %v, err=%v", topic, err))
	}

	settings := spec.Settings()
	err = settings.SetLogLevel(levelStr)
	if err != nil {
		return err
	}

	if pipeline != nil {
		if pipeline.RuntimeContext() == nil {
			return errors.New("Pipeline doesn't have the runtime context")
		}
		if pipeline.RuntimeContext().Service(base.PIPELINE_SUPERVISOR_SVC) == nil {
			return errors.New("Pipeline doesn't have the PipelineSupervisor registered")
		}
		supervisor := pipeline.RuntimeContext().Service(base.PIPELINE_SUPERVISOR_SVC).(*pipeline_svc.PipelineSupervisor)

		if supervisor != nil {
			err := supervisor.SetPipelineLogLevel(levelStr)
			return err
		}
	}
	return nil
}
