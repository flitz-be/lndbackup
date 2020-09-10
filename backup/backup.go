package backup

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	_ "gocloud.dev/blob/gcsblob"
)

// ChannelSnapshot uploads a channel backup snapshot muiltchannel
// blob to a blob storage bucket
func ChannelSnapshot(ctx context.Context, bucketURL string, backup []byte) error {
	bucket, err := OpenBucket(ctx, bucketURL)
	if err != nil {
		log.Printf("Failed to open bucket: %v", err)
		return err
	}
	defer bucket.Close()

	backupFileName := getName()

	// Open a *blob.Writer for the blob at blobKey.
	writer, err := bucket.NewWriter(ctx, backupFileName, nil)
	if err != nil {
		log.Printf("Failed to write %q: %v", backupFileName, err)
		return err
	}
	defer writer.Close()

	// Copy the data.
	log.Printf("Backup started: %v", backupFileName)
	_, err = io.Copy(writer, bytes.NewReader(backup))
	if err != nil {
		log.Printf("Failed to copy data: %v", err)
		return err
	}
	log.Printf("Backup finished: %v", backupFileName)

	return nil
}

func getName() string {
	now := time.Now()
	nsec := now.Local().UnixNano()
	return fmt.Sprintf("channel-backup-%v", (nsec / 1000 / 1000))
}
