package provider

import "sync"

// atmosMu serializes all calls into the Atmos library.
//
// The Atmos library uses package-level mutable state (e.g., mergedConfigFiles in
// pkg/config/load.go) that is explicitly documented as not safe for concurrent use.
// Terraform invokes ReadDataSource concurrently for independent data sources, so
// without this mutex, concurrent calls corrupt shared state and trigger os.Exit(1)
// via CheckErrorPrintAndExit, killing the gRPC plugin process.
var atmosMu sync.Mutex
