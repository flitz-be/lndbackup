package backup

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
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

func bucketPath(path string) string {
	path = strings.TrimLeft(path, "/")
	path = filepath.Clean(path) + "/"
	return path
}
