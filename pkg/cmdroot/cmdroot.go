package cmdroot

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	programName string
	programPath string
	cmdRoot     *cobra.Command
	options     initOptions
)

func init() {
	// get binary location
	file, _ := exec.LookPath(os.Args[0])
	programPath, programName = filepath.Split(file)

	cmdRoot = getRootCommand()
	cobra.OnInitialize(initConfig)
}

type initOptions struct {
	report  bool
	monitor bool
}

// InitOption for cmdroot
type InitOption func(o *initOptions)

// WithReport enable report service AdminAddress/MonitorAddress etc.
// It will read service info from config:
//
//    AdminAddress : admin address of service
//    MonitorAddress : monitor address for prometheus
//    LocEndpoints : etcd address
func WithReport() InitOption {
	return func(o *initOptions) {
		o.report = true
	}
}

// WithMonitor enable prometheus report
// It will read `MonitorAddress` for metrics bind address
func WithMonitor() InitOption {
	return func(o *initOptions) {
		o.monitor = true
	}
}

// InitCommand set command description and logger
//
// Input:
//    short: short description
//    long:  longer description
func InitCommand(short, long string, opts ...InitOption) {
	cmdRoot.Short = short
	cmdRoot.Long = long
	for _, opt := range opts {
		opt(&options)
	}
}

func initConfig() {
}

// AddCommand add child command,
func AddCommand(cmd *cobra.Command) {
	cmdRoot.AddCommand(cmd)
}

// Execute command
func Execute() {
	rand.Seed(time.Now().UnixNano())
	if c, err := cmdRoot.ExecuteC(); err != nil {
		if isUserError(err) {
			cmdRoot.Println("")
			cmdRoot.Println(c.UsageString())
		}

		os.Exit(-1)
	}
}
