# Description

This repository is the companion of [grafana-backup-runner](https://github.com/invincibleinfra/grafana-backup-runner).
It restores from backups generated by the preceeding tool

# Setup

To build the executable, run `make build`.

# Usage

If you have a shared AWS configuration residing in `~/.aws`, then invoke the executable like so:

```
./grafana-backup-restore -grafanaURL ${GRAFANA_URL} -s3Bucket ${S3_BUCKET_NAME} -useSharedConfig -backupPath ${BACKUP_PATH_IN_S3}
```

If you are not using a shared AWS config, then omit the `-useSharedConfig` and set all AWS credential environment variables
described [here](https://github.com/invincibleinfra/grafana-backup-runner).
