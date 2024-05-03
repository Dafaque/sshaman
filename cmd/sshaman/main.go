package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Dafaque/sshaman/internal/config"
	"github.com/Dafaque/sshaman/internal/credentials"
	"github.com/Dafaque/sshaman/internal/credentials/local"
)

const (
	emptyString string = ""
)

var (
	flagLocal  bool
	flagRemote bool

	flagAlias string

	//? MARK: - list options
	flagList bool

	//? MARK: - Add connection
	flagHost           string
	flagPort           int
	flagUser           string
	flagKeyFilePath    string
	flagSkipPassword   bool
	flagSkipPassphrase bool

	//! MARK: - Danger zone
	flagForce bool
	flagDrop  bool
)

const (
	cmdAdd     string = "add"
	cmdConnect string = "connect"
	cmdList    string = "list"
	cmdDelete  string = "delete"
	cmdDrop    string = "drop"
)

type operation func(credentials.Manager, credentials.Manager) error

var errRemoteNotConfigured error = errors.New("remote credentials manager is not configured")

func main() {
	// MARK: - Add connection
	addFlags := flag.NewFlagSet(cmdAdd, flag.ExitOnError)
	addFlags.StringVar(&flagAlias, "alias", emptyString, "new ssh connection's alias")
	addFlags.StringVar(&flagHost, "host", emptyString, "new ssh connection's address")
	addFlags.IntVar(&flagPort, "port", 22, "new ssh connection's port")
	addFlags.StringVar(&flagUser, "user", emptyString, "new ssh connection's user")
	addFlags.StringVar(&flagKeyFilePath, "key", emptyString, "new ssh connection's key file path")
	addFlags.BoolVar(&flagSkipPassword, "skip-password", false, "skip password prompt")
	addFlags.BoolVar(&flagSkipPassphrase, "skip-passphrase", false, "skip key's passphrase prompt")
	addFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	addFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplementer)") //@todo implement
	addFlags.BoolVar(&flagForce, "force", false, "force operation")

	// MARK: - Connect
	conectFlags := flag.NewFlagSet(cmdConnect, flag.ExitOnError)
	conectFlags.StringVar(&flagAlias, "alias", emptyString, "ssh connection's allias to conect")
	conectFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	conectFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement

	// MARK: - List connections
	listFlags := flag.NewFlagSet(cmdList, flag.ExitOnError)
	listFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	listFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement

	// MARK: - Delete connection
	delFlags := flag.NewFlagSet(cmdDelete, flag.ExitOnError)
	delFlags.StringVar(&flagAlias, "alias", emptyString, "ssh connection's allias to conect")
	delFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	delFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement

	// MARK: - Drop all connections
	dropFlags := flag.NewFlagSet(cmdDrop, flag.ExitOnError)
	dropFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	dropFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement
	dropFlags.BoolVar(&flagForce, "force", false, "force operation")

	// MARK: - App details
	ver := flag.Bool("version", false, "show app details")
	flag.Parse()

	if *ver {
		i, _ := debug.ReadBuildInfo()
		fmt.Println("goVersion:", i.GoVersion)
		fmt.Println("module:", i.Main.Path)
		fmt.Println("version:", i.Main.Version)
		fmt.Println("sum:", i.Main.Sum)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		addFlags.Usage()
		conectFlags.Usage()
		listFlags.Usage()
		delFlags.Usage()
		dropFlags.Usage()
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	localManager, err := local.NewManager(cfg)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := localManager.Done(); err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}()

	var flagset *flag.FlagSet
	var op operation
	switch os.Args[1] {
	case cmdAdd:
		flagset = addFlags
		op = addConnection
	case cmdConnect:
		flagset = conectFlags
		op = connect
	case cmdList:
		flagset = listFlags
		op = listCredentials
	case cmdDelete:
		flagset = delFlags
		op = deleteCredentials
	case cmdDrop:
		flagset = dropFlags
		op = dropCredentials
	default:
		addFlags.Usage()
		conectFlags.Usage()
		listFlags.Usage()
		delFlags.Usage()
		dropFlags.Usage()
		flag.Usage()
		os.Exit(1)
	}

	errParse := flagset.Parse(os.Args[2:])
	if err != nil {
		println(errParse.Error())
		flag.Usage()
		os.Exit(1)
	}

	if err := op(localManager, nil); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
