package cmd
import (
"net"
 "fmt"
"flag"
 "os"
"strconv"
"github.com/openebs/mayaserver/lib/config"

	"github.com/mitchellh/cli"
	"github.com/openebs/mayaserver/lib/flaghelper"
)
type StatusCommand struct {
	Meta
	Ui   cli.Ui
	args []string
}
func verify () {
  flag.Usage = usage
  flag.Parse()

  args := flag.Args()
  if len(args) < 1 {
      fmt.Fprintf(os.Stderr, "Input port is missing.")
      os.Exit(1)
  }

  port := args[0]
  _, err := strconv.ParseUint(port, 10, 16)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Invalid port %q: %s\n", port, err)
    os.Exit(1)
  }

  ln, err := net.Listen("tcp", ":" + port)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Can't listen on port %q: %s\n", port, err)
    os.Exit(1)
  }

  err = ln.Close()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Couldn't stop listening on port %q: %s\n", port, err)
    os.Exit(1)
  }

  fmt.Printf("TCP Port %q is available\n", port)
  os.Exit(0)
}
func (c *StatusCommand) readMayaConfig() *config.MayaConfig {
	var configPath []string
	cmdConfig := &config.MayaConfig{
		Ports: &config.Ports{},
	}

	flags := flag.NewFlagSet("up", flag.ContinueOnError)
	flags.Usage = func() { c.Ui.Error(c.Help()) }

	flags.Var((*flaghelper.StringFlag)(&configPath), "config", "config")
	flags.StringVar(&cmdConfig.BindAddr, "bind", "", "")
	flags.StringVar(&cmdConfig.NodeName, "name", "m-apiserver", "")

	flags.StringVar(&cmdConfig.DataDir, "data-dir", "", "")
	flags.StringVar(&cmdConfig.LogLevel, "log-level", "", "")
}

func (c *StatusCommand) Run(args []string) int {
  if config.vmbox_status == 'running'
 {

  mconfig := c.readMayaConfig() 

     c.Ui.Output(out)
     out :=     fmt.Sprintf("Name          IP            Ports       Status    \n%-16s%-16s%d\t%-16s",
			mconfig.NodeName,
			mconfig.BindAddr,
			mconfig.Ports,
			Status)
    
    }
     }
