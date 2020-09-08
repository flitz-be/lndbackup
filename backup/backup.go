package backup

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	_ "gocloud.dev/blob/gcsblob"
)

// ChannelSnapshot uploads a channel backup snapshot muiltchannel
// blob to a blob storage bucket
func ChannelSnapshot(ctx context.Context, bucketURL string, snapshot lnrpc.ChanBackupSnapshot) {
	bucket, err := OpenBucket(ctx, bucketURL)

	if err != nil {
		log.Printf("Failed to open bucket: %v\n", err)
	}
	defer bucket.Close()

	backupFileName := getName()

	// Open a *blob.Writer for the blob at blobKey.
	writer, err := bucket.NewWriter(ctx, backupFileName, nil)
	if err != nil {
		log.Printf("Failed to write %q: %v\n", backupFileName, err)
	}
	defer writer.Close()

	// Copy the data.
	// bytes.NewReader(backup.MultiChanBackup.MultiChanBackup
	copied, err := io.Copy(writer, bytes.NewReader(snapshot.MultiChanBackup.MultiChanBackup))
	log.Printf("filename: %v, lenght: %v\n", backupFileName, copied)
	if err != nil {
		log.Printf("Failed to copy data: %v\n", err)
	}
}

func getName() string {
	now := time.Now()
	nsec := now.Local().UnixNano()
	return fmt.Sprintf("channel-backup-%v", (nsec / 1000 / 1000))
}
