package backup

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/s3blob"
	"gocloud.dev/gcp"
)

// OpenBucket returns a bucket using a bucket url.
// It will prefix any object stored in it.
func OpenBucket(ctx context.Context, bucketName string) (*blob.Bucket, error) {
	url, err := url.Parse(bucketName)
	if err != nil {
		return nil, fmt.Errorf("Bucket url invalid: %v", err)
	}

	var bucket *blob.Bucket
	switch url.Scheme {
	case "gs":
		bucket, err = gsBucket(ctx, url.Hostname())
	case "s3":
		bucket, err = s3Bucket(ctx, url.Hostname())
	default:
		return nil, fmt.Errorf("Bucket provider invalid: invalid provider %s", url.Scheme)
	}

	bucket = blob.PrefixedBucket(bucket, bucketPath(url.Path))
	return bucket, err
}

func gsBucket(ctx context.Context, name string) (*blob.Bucket, error) {
	// See here for information on credentials:
	// https://cloud.google.com/docs/authentication/getting-started
	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}
	c, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
	if err != nil {
		return nil, err
	}
	return gcsblob.OpenBucket(ctx, c, name, nil)
}

func s3Bucket(ctx context.Context, name string) (*blob.Bucket, error) {
	//Based on:
	// - https://gocloud.dev/howto/blob/#s3-compatible
	// - https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/
	key := os.Getenv("S3_KEY")
	secret := os.Getenv("S3_SECRET")
	endpoint := os.Getenv("S3_ENDPOINT")
	region := os.Getenv("S3_REGION")
	if key == "" || secret == "" || region == "" {
		return nil, fmt.Errorf("key, secret and region should be set: %s, %s, %s", key, secret, region)
	}
	var ptrEndpoint *string
	if endpoint != "" {
		ptrEndpoint = aws.String(endpoint)
	}
	ptrRegion := aws.String(region)
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    ptrEndpoint,
		Region:      ptrRegion,
	})
	if err != nil {
		return nil, err
	}
	return s3blob.OpenBucket(ctx, sess, name, nil)
}

func bucketPath(path string) string {
	path = strings.TrimLeft(path, "/")
	path = filepath.Clean(path) + "/"
	return path
}
