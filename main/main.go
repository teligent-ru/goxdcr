package main

import (
	"flag"
	"fmt"
	"os"

	ap "github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/adminport"
	rm "github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/replication_manager"
	c "github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/mock_services"
)

var done = make(chan bool)

var options struct {
	sourceClusterAddr      string //source cluster addr
	username        string //username on source cluster
	password        string //password on source cluster	
	// temp params for mock services. should be removed later
	sourceBucket          string // source bucket
	targetBucket          string // target bucket
	numConnPerKV  int    // number of nozzles per source node
	numOutgoingConn  int    // number of nozzles per target node
}

func argParse() {
	flag.StringVar(&options.sourceClusterAddr, "sourceClusterAddr", "127.0.0.1:9000",
		"source cluster address")
	flag.StringVar(&options.username, "username", "Administrator", "username to cluster admin console")
	flag.StringVar(&options.password, "password", "welcome", "password to Cluster admin console")
		
	flag.StringVar(&options.sourceBucket, "sourceBucket", "default",
		"bucket to replicate from")
	flag.StringVar(&options.targetBucket, "targetBucket", "target",
		"bucket to replicate to")
	flag.IntVar(&options.numConnPerKV, "numConnPerKV", 2,
		"number of nozzles per source node")
	flag.IntVar(&options.numOutgoingConn, "numOutgoingConn", 3,
		"number of nozzles per target node")
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage : %s [OPTIONS] <sourceClusterAddr> \n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	argParse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	options.sourceClusterAddr = args[0]
	
	c.SetTestOptions(options.sourceBucket, options.targetBucket, options.sourceClusterAddr, options.sourceClusterAddr, options.username, options.password, options.numConnPerKV, options.numOutgoingConn)
	xdcrTopologyService := new(c.MockXDCRTopologySvc)
	hostAddr, err := xdcrTopologyService.MyHost()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting host address \n")
		os.Exit(1)
	}
	rm.Initialize(new(c.MockMetadataSvc), new(c.MockClusterInfoSvc), xdcrTopologyService, new(c.MockReplicationSettingsSvc))
	go ap.MainAdminPort(hostAddr)
	<-done
}
