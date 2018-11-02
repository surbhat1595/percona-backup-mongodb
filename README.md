# MongoDB backup tool
[![codecov](https://codecov.io/gh/percona/mongodb-backup/branch/master/graph/badge.svg?token=TiuOmTfp2p)](https://codecov.io/gh/percona/mongodb-backup)

Progress:
- [x] Oplog tailer
- [x] Oplog applier
- [x] S3 streamer
- [x] Mongodump backup
- [x] Mongodump restore
- [x] Agent selection
- [x] Replica Set Backup
- [x] Sharded Cluster Backup

## Building

Building the project requires:
1. Go 1.11 or above
1. make
1. upx *(optional)*

To build the project *(from the project dir)*:
```
make
```

A successful build outputs binaries: 
1. **pmb-admin**: A command-line interface for controlling the backup system
1. **pmb-agent**: An agent that executes backup/restore actions on a database host
1. **pmb-coordinator**: A server that coordinates backup system actions

## Testing

The integration testing launches a MongoDB cluster in Docker containers. *'docker'* and *'docker-compose'* is required.

To run the tests *(may require 'sudo')*:
```
make test-full
```

To tear-down the test *(and containers, data, etc)*:
```
make test-full-clean
```

## Run in Docker

### Build Docker images

To build the Docker images:
    ```
    make docker-build
    ```

### Coordinator
*Note: '/data/percona-mongodb-backup' must be owned by UID 100*
    ```
    docker run -d --rm --name mongodb-backup-coordinator \
        -p 10000:10000 \
        -p 10001:10001 \
        -v /data/percona-mongodb-backup:/data/percona-mongodb-backup \
    mongodb-backup-coordinator
    ```
