package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Dafaque/sshaman/v2/internal/config"
	"github.com/Dafaque/sshaman/v2/internal/credentials"
)

const (
	emptyString string = ""
)

var (
	flagName string

	//? MARK: - list options
	flagList bool

	//? MARK: - Add connection
	flagHost           string
	flagPort           int
	flagUser           string
	flagKeyFilePath    string
	flagSkipPassword   bool
	flagSkipPassphrase bool

	//? MARK: - Import
	flagDryRun bool

	//! MARK: - Danger zone
	flagForce bool
)

const (
	cmdAdd     string = "add"
	cmdConnect string = "connect"
	cmdList    string = "list"
	cmdRmove   string = "remove"
	cmdDrop    string = "drop"
	cmdExport  string = "export"
	cmdImport  string = "import"
)

type operation func(*credentials.Manager) error

var errRemoteNotConfigured error = errors.New("remote credentials manager is not configured")

func main() {
	// MARK: - Add connection
	addFlags := flag.NewFlagSet(cmdAdd, flag.ExitOnError)
	addFlags.StringVar(&flagName, "n", emptyString, "new ssh connection's alias")
	addFlags.StringVar(&flagHost, "h", emptyString, "new ssh connection's address")
	addFlags.IntVar(&flagPort, "p", 22, "new ssh connection's port")
	addFlags.StringVar(&flagUser, "u", emptyString, "new ssh connection's user")
	addFlags.StringVar(&flagKeyFilePath, "k", emptyString, "new ssh connection's key file path")
	addFlags.BoolVar(&flagSkipPassword, "no-pw", false, "skip password prompt")
	addFlags.BoolVar(&flagSkipPassphrase, "no-pp", false, "skip key's passphrase prompt")
	addFlags.BoolVar(&flagForce, "f", false, "force operation")

	// MARK: - Connect
	conectFlags := flag.NewFlagSet(cmdConnect, flag.ExitOnError)
	conectFlags.StringVar(&flagName, "n", emptyString, "ssh connection's allias to conect")

	// MARK: - List connections
	listFlags := flag.NewFlagSet(cmdList, flag.ExitOnError)

	// MARK: - Delete connection
	rmFlags := flag.NewFlagSet(cmdRmove, flag.ExitOnError)
	rmFlags.StringVar(&flagName, "n", emptyString, "ssh connection's allias to conect")

	// MARK: - Drop all connections
	dropFlags := flag.NewFlagSet(cmdDrop, flag.ExitOnError)
	dropFlags.BoolVar(&flagForce, "f", false, "force operation")

	// MARK: - Export
	exportFlags := flag.NewFlagSet(cmdExport, flag.ExitOnError)
	exportFlags.BoolVar(&flagSkipPassword, "no-pw", false, "skip password prompt")

	// MARK: - Import
	importFlags := flag.NewFlagSet(cmdImport, flag.ExitOnError)
	importFlags.BoolVar(&flagDryRun, "dry", false, "view what would be imported")
	importFlags.BoolVar(&flagSkipPassword, "no-pw", false, "skip password prompt")

	// MARK: - App details
	ver := flag.Bool("v", false, "show app details")
	flag.Parse()

	usage := func() {
		addFlags.Usage()
		conectFlags.Usage()
		// listFlags.Usage() //? no args
		rmFlags.Usage()
		dropFlags.Usage()
		importFlags.Usage()
		exportFlags.Usage()
		flag.Usage()
		os.Exit(1)
	}

	if *ver {
		i, _ := debug.ReadBuildInfo()
		fmt.Println("goVersion:", i.GoVersion)
		fmt.Println("module:", i.Main.Path)
		fmt.Println("version:", i.Main.Version)
		fmt.Println("sum:", i.Main.Sum)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		usage()
	}

	cfg, err := config.NewConfig()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	manager, err := credentials.New(cfg)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := manager.Done(); err != nil {
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
	case cmdRmove:
		flagset = rmFlags
		op = removeCredentials
	case cmdDrop:
		flagset = dropFlags
		op = dropCredentials
	case cmdExport:
		flagset = exportFlags
		op = exportCredentials
	case cmdImport:
		flagset = importFlags
		op = importCredentials
	default:
		usage()
	}

	errParse := flagset.Parse(os.Args[2:])
	if err != nil {
		println(errParse.Error())
		flag.Usage()
		os.Exit(1)
	}

	if err := op(manager); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
