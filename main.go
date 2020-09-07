package main

import (
	"context"
	"flag"
	"log"
	"os"

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

	// maxMsgRecvSize is the largest message our client will receive. We
	// set this to ~50Mb atm.
	maxMsgRecvSize = grpc.MaxCallRecvMsgSize(1 * 1024 * 1024 * 50)

	// Defaults values
	defaultRPCHost     = getEnv("RPC_HOST", "localhost")
	defaultRPCPort     = getEnv("RPC_PORT", "10009")
	defaultTLSCertPath = getEnv("TLS_CERT_PATH", "/root/.lnd")
	defaultMacaroonDir = getEnv("MACAROON_PATH", "")
	defaultNetwork     = getEnv("NETWORK", "mainnet")

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
)

func main() {
	flag.Parse()

	log.Printf("Starting Lightning Static Channel Backup Version=%v GitCommit=%v", version, gitCommit)

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

	backupUpdates, _, _ := client.Client.SubscribeChannelBackups(context.Background())

	for {
		select {
		case backups := <-backupUpdates:
			log.Printf(backups.MultiChanBackup.String())
		case <-context.Background().Done():
		}
	}
}
