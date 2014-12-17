// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package base

import (
	"errors"
)

//constants
var DefaultConnectionSize = 10
var DefaultPoolName = "default"

var LocalHostName = "localhost"

// URL Paths for retrieving cluster info
var PoolsPath = "/pools"
var NodesPath = "/pools/nodes"
var SSLPortsPath = "/nodes/self/xdcrSSLPorts"

// keys used in parsing cluster info
var NodesKey = "nodes"
var HostNameKey = "hostname"
var ThisNodeKey = "thisNode"
var SSLPortKey = "httpsMgmt" 

// URL related constants
var UrlDelimiter = "/"
var UrlPortNumberDelimiter = ":"

// http request method types
const (
	MethodGet = "GET"
	MethodPost = "POST"
	MethodDelete = "DELETE"
)

// delimiter for multiple parts in a key
var KeyPartsDelimiter = "_"

//constants for adminport
var AdminportUrlPrefix = UrlDelimiter

// used as default value for tests
var AdminportNumber uint16 = 13000
var GometaRequestPortNumber uint16 = 11000

// AdminportReadTimeout timeout, in milliseconds, is read timeout for
// golib's http server.
var AdminportReadTimeout = 0
// AdminportWriteTimeout timeout, in milliseconds, is write timeout for
// golib's http server.
var AdminportWriteTimeout = 0

//outgoing nozzle type
type XDCROutgoingNozzleType int

const (
	Xmem XDCROutgoingNozzleType = iota
	Capi XDCROutgoingNozzleType = iota
)

const (
	PIPELINE_SUPERVISOR_SVC string = "PipelineSupervisor"
	CHECKPOINT_MGR_SVC string = "CheckpointManager"
	STATISTICS_MGR_SVC string = "StatisticsManager"
)

// supervisor related constants
const(
	ReplicationManagerSupervisorId = "ReplicationManagerSupervisor"
	PipelineMasterSupervisorId = "PipelineMasterSupervisor"
	AdminportSupervisorId = "AdminportSupervisor"
	PipelineSupervisorIdPrefix = "PipelineSupervisor_"
)

// constants for integer parsing
var ParseIntBase    = 10
var ParseIntBitSize = 64

var ErrorNotMyVbucket = errors.New("NOT_MY_VBUCKET")

// constants used for remote cluster references
const (
	RemoteClustersPath  = "pools/default/remoteClusters"
	
	RemoteClusterUuid   = "uuid"
	RemoteClusterName  = "name"
	RemoteClusterHostName = "hostname"
	RemoteClusterUserName = "username"
	RemoteClusterPassword = "password"
	RemoteClusterDemandEncryption = "demandEncryption"
	RemoteClusterCertificate = "certificate"
	RemoteClusterUri = "uri"
	RemoteClusterValidateUri = "validateURI"
	RemoteClusterDeleted = "deleted"
)

// constant used by more than one rest apis
const (
 	JustValidate = "just_validate"
)

