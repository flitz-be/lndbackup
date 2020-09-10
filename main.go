package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/budacom/lnd-backup/backup"
	"github.com/lightninglabs/lndclient"

	"google.golang.org/grpc"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

var (
	// Set during go build
	version   string
	gitCommit string

	signalsChannel = make(chan os.Signal, 1)
	quit           = make(chan struct{})

	backupsPending = false

	// maxMsgRecvSize is the largest message our client will receive. We
	// set this to ~50Mb atm.
	maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 50)

	// Defaults values
	defaultRPCHost     = getEnv("RPC_HOST", "localhost")
	defaultRPCPort     = getEnv("RPC_PORT", "10009")
	defaultTLSCertPath = getEnv("TLS_CERT_PATH", "/root/.lnd")
	defaultMacaroonDir = getEnv("MACAROON_PATH", "")
	defaultNetwork     = getEnv("NETWORK", "mainnet")
	defaultBucketURL   = getEnv("BUCKET_URL", "")

	// Command-line flags
	rpcHost = flag.String("rpc-host", defaultRPCHost,
		"Lightning node RPC host. The default value can be overwritten by RPC_HOST environment variable.")
	rpcPort = flag.String("rpc-port", defaultRPCPort,
		"Lightning node RPC port. The default value can be overwritten by RPC_PORT environment variable.")
	tlsCertPath = flag.String("tls-cert-path", defaultTLSCertPath,
		"The path to the tls certificate. The default value can be overwritten by TLS_CERT_PATH environment variable.")
	macaroonDir = flag.String("macaroon-dir", defaultMacaroonDir,
		"The directory where macaroons are stored. The default value can be overwritten by MACAROON_DIR environment variable.")
	network = flag.String("network", defaultNetwork,
		"The chain network to operate on. The default value can be overwritten by NETWORK environment variable.")
	bucketURL = flag.String("bucket-url", defaultBucketURL,
		"The bucket url to backup the snapshot. The default value can be overwritten by BUCKET_URL environment variable.")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Starting Lightning Static Channel Backup Version=%v GitCommit=%v", version, gitCommit)

	handleSignals()

	client, err := lndclient.NewLndServices(&lndclient.LndServicesConfig{
		LndAddress:  *rpcHost,
		Network:     lndclient.Network(*network),
		MacaroonDir: *macaroonDir,
		TLSPath:     *tlsCertPath,
		// Use the default lnd version check which checks for version
		// v0.11.0 and requires all build tags.
		CheckVersion: nil,
	})
	if err != nil {
		log.Printf("cannot connect to lightning services: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	log.Printf("Initial backup")
	snapShot, err := client.Client.ChannelBackups(ctx)
	err = backup.ChannelSnapshot(ctx, *bucketURL, snapShot)
	if err != nil {
		os.Exit(1)
	}

	backupUpdates, _, _ := client.Client.SubscribeChannelBackups(ctx)
	log.Printf("Subscribed to channel backups")

ReadBackupUpdates:
	for {
		select {
		case channelSnapshot := <-backupUpdates:
			backupsPending = true
			multiChanBackup := channelSnapshot.GetMultiChanBackup()
			backup.ChannelSnapshot(ctx, *bucketURL, multiChanBackup.MultiChanBackup)
			backupsPending = false
		case <-quit:
			break ReadBackupUpdates
		}
	}
}

func handleSignals() {
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-signalsChannel:
				if backupsPending == true {
					log.Println("Waiting for pending backups")
				}
				log.Println("Shutting Down Gracefully")
				close(quit)
				return
			}
		}
	}()
}
