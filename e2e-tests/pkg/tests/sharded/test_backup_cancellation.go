package sharded

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go"
	"gopkg.in/yaml.v2"

	"github.com/percona/percona-backup-mongodb/pbm"
	pbmt "github.com/percona/percona-backup-mongodb/pbm"
)

func (c *Cluster) BackupCancellation(storage string) {
	bcpName := c.Backup()
	ts := time.Now()
	time.Sleep(1 * time.Second)
	c.printBcpList()
	log.Println("canceling backup", bcpName)
	o, err := c.pbm.RunCmd("pbm", "cancel-backup")
	if err != nil {
		log.Fatalf("Error: cancel backup '%s'.\nOutput: %s\nStderr:%v", bcpName, o, err)
	}

	time.Sleep(3 * time.Second)

	checkNoBackupFiles(bcpName, storage)

	log.Println("check backup state")
	m, err := c.mongopbm.GetBackupMeta(bcpName)
	if err != nil {
		log.Fatalf("Error: get metadata for backup %s: %v", bcpName, err)
	}

	if m.Status != pbm.StatusCancelled {
		log.Fatalf("Error: wrong backup status, expect %s, got %v", pbm.StatusCancelled, m.Status)
	}

	needToWait := pbmt.WaitActionStart + time.Second - time.Since(ts)
	if needToWait > 0 {
		log.Printf("waiting for the lock to be released for %s", needToWait)
		time.Sleep(needToWait)
	}
}

func checkNoBackupFiles(backupName, conf string) {
	log.Println("check no artefacts left for backup", backupName)
	buf, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatalln("Error: unable to read config file:", err)
	}

	var cfg pbm.Config
	err = yaml.UnmarshalStrict(buf, &cfg)
	if err != nil {
		log.Fatalln("Error: unmarshal yaml:", err)
	}

	stg := cfg.Storage

	endopintURL := awsurl
	if stg.S3.EndpointURL != "" {
		eu, err := url.Parse(stg.S3.EndpointURL)
		if err != nil {
			log.Fatalln("Error: parse EndpointURL:", err)
		}
		endopintURL = eu.Host
	}

	mc, err := minio.NewWithRegion(endopintURL, stg.S3.Credentials.AccessKeyID, stg.S3.Credentials.SecretAccessKey, false, stg.S3.Region)
	if err != nil {
		log.Fatalln("Error: NewWithRegion:", err)
	}

	for object := range mc.ListObjects(stg.S3.Bucket, stg.S3.Prefix, true, nil) {
		if object.Err != nil {
			fmt.Println("Error: ListObjects: ", object.Err)
			continue
		}

		if strings.Contains(object.Key, backupName) {
			log.Fatalln("Error: failed to delete lefover", object.Key)
		}
	}
}