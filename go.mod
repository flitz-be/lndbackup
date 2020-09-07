module github.com/budacom/lnd-backup

go 1.15

replace github.com/lightninglabs/lndclient => ../../lndclient

require (
	github.com/btcsuite/btcd v0.20.1-beta.0.20200730232343-1db1b6f8217f
	github.com/gogo/protobuf v1.1.1
	github.com/lightninglabs/lndclient v1.0.0
	github.com/lightningnetwork/lnd v0.11.0-beta
	github.com/prometheus/client_golang v0.9.3
	google.golang.org/grpc v1.24.0
	gopkg.in/macaroon.v2 v2.1.0
)
