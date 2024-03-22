package utils

import "os"

var (
	Env                    = os.Getenv("ENV")
	Env_TracingServiceName = os.Getenv("TRACING_SERVICE_NAME")
	Env_OLTPEndpoint       = os.Getenv("OLTP_ENDPOINT")

	// Read-write user on the local CH node, can be admin user
	CH_WRITE_DSN = os.Getenv("CH_WRITE_DSN")

	// Read only user on the local CH node
	CH_READ_DSN = os.Getenv("CH_READ_DSN")

	// DNS host for requests to the leader. E.g. "leader-rafthouse.svc.default.cluster.local"
	LEADER_HOST = os.Getenv("LEADER_HOST")

	DB_PATH = GetEnvOrDefault("DB_PATH", "/etc/rafthouse")
)
