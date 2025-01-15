package couchbase

import (
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

var (
	Cluster *gocb.Cluster
	Bucket  *gocb.Bucket
)

// InitCouchbase initializes the connection to the Couchbase cluster and bucket
func InitCouchbase(username, password, bucketName string) {
	var err error

	// Connect to the Couchbase Cluster
	Cluster, err = gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Couchbase: %v", err)
	}

	// Connect to the specified Bucket
	Bucket = Cluster.Bucket(bucketName)
	if err = Bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		log.Fatalf("Failed to connect to bucket: %v", err)
	}

	log.Println("Connected to Couchbase successfully!")
}

// GetCollection returns the default collection for performing document operations
// This can be used for CRUD operations
func GetCollection() *gocb.Collection {
	return Bucket.DefaultCollection()
}
