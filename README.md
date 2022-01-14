# LND Static Channels Backup
Streams SCB's from LND to a object storage provider.
Supports GCS and S3-compatible providers.
Tested with Digital Ocean Spaces.

## GCS
GCP credentials should be in default location.
## S3
Code based on [gocloud docs](https://gocloud.dev/howto/blob/#s3-compatible) and [Digital Ocean S3 compatibility docs](https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/).

```
BUCKET_URL=s3://bucket_name/sub/folder
S3_KEY=Access_Key_ID
S3_SECRET=Access_Key_Secret
S3_ENDPOINT=endpoint (leave blank if using AWS)
S3_REGION=s3_region (only used when provider is AWS)
```