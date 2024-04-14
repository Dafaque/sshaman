package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"golang.org/x/term"

	"github.com/Dafaque/sshaman/internal/config"
	"github.com/Dafaque/sshaman/internal/credentials"
	"github.com/Dafaque/sshaman/internal/credentials/local"
	"github.com/Dafaque/sshaman/pkg/clients/ssh"
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

func main() {
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

	conectFlags := flag.NewFlagSet(cmdConnect, flag.ExitOnError)
	conectFlags.StringVar(&flagAlias, "alias", emptyString, "ssh connection's allias to conect")
	conectFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	conectFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement

	listFlags := flag.NewFlagSet(cmdList, flag.ExitOnError)
	listFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	listFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement

	delFlags := flag.NewFlagSet(cmdDelete, flag.ExitOnError)
	delFlags.StringVar(&flagAlias, "alias", emptyString, "ssh connection's allias to conect")
	delFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	delFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement

	dropFlags := flag.NewFlagSet(cmdDrop, flag.ExitOnError)
	dropFlags.BoolVar(&flagLocal, "local", true, "use local storage")
	dropFlags.BoolVar(&flagRemote, "remote", false, "use remote storage (unimplemented)") //@todo implement
	dropFlags.BoolVar(&flagForce, "force", false, "force operation")

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

var errRemoteNotConfigured error = errors.New("remote credentials manager is not configured")

func connect(local credentials.Manager, remote credentials.Manager) error {
	var creds *credentials.Credentials
	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		remoteCreds, err := remote.Get(flagAlias)
		if err != nil {
			return err
		}
		creds = remoteCreds
	}
	if flagLocal && creds == nil {
		localCreds, err := local.Get(flagAlias)
		if err != nil {
			return err
		}
		creds = localCreds
	}
	if creds == nil {
		return errors.New("no credentials source given")
	}
	cl, err := ssh.NewSshClient(creds)
	if err != nil {
		return err
	}
	return cl.Loop()
}

func addConnection(local credentials.Manager, remote credentials.Manager) error {
	creds, err := makeNewCredentials()
	if err != nil {
		return err
	}

	if flagLocal {
		if err := local.Set(flagAlias, creds, flagForce); err != nil {
			return err
		}
		fmt.Println("local credentials added for", flagAlias)
	}

	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		if err := remote.Set(flagAlias, creds, flagForce); err != nil {
			return err
		}
		fmt.Println("remote credentials added for", flagAlias)
	}

	return nil
}

func makeNewCredentials() (*credentials.Credentials, error) {
	if flagHost == emptyString {
		return nil, errors.New("host required")
	}

	if flagUser == emptyString {
		return nil, errors.New("user required")
	}

	if flagAlias == emptyString {
		return nil, errors.New("alias required")
	}

	var creds credentials.Credentials = credentials.Credentials{
		Alias:    flagAlias,
		Host:     flagHost,
		Port:     flagPort,
		Username: flagUser,
	}

	fd := int(os.Stdin.Fd())
	var password string
	var passphrase []byte
	var key []byte
	if flagKeyFilePath != emptyString {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		if flagKeyFilePath == "~" {
			// In case of "~", which won't be caught by the "else if"
			flagKeyFilePath = home
		} else if strings.HasPrefix(flagKeyFilePath, "~/") {
			flagKeyFilePath = path.Join(home, flagKeyFilePath[2:])
		}
		file, err := os.ReadFile(flagKeyFilePath)
		if err != nil {
			return nil, err
		}
		key = file
		if !flagSkipPassphrase {
			fmt.Printf("Enter %s key passphrase: ", flagKeyFilePath)
			pp, err := term.ReadPassword(fd)
			println()
			if err != nil {
				return nil, err
			}
			passphrase = pp
		}
	}
	if !flagSkipPassword {
		fmt.Printf("Enter %s's password for %s: ", flagUser, flagHost)
		pass, err := term.ReadPassword(fd)
		println()
		if err != nil {
			return nil, err
		}
		password = string(pass)
	}

	if len(password) > 0 {
		creds.Password = &password
	}
	if len(passphrase) > 0 {
		creds.Passphrase = passphrase
	}
	if len(key) > 0 {
		creds.Key = key
	}

	return &creds, nil
}

func dropCredentials(local credentials.Manager, remote credentials.Manager) error {
	if !flagForce {
		return errors.New("this operation will delete all your data. if you are sure of what you are doing, use the flag -force")
	}
	if flagLocal {
		if err := local.Drop(); err != nil {
			return err
		}
		fmt.Println("local credentials cleared")
	}
	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		if err := remote.Drop(); err != nil {
			return err
		}
		fmt.Println("remote credentials cleared")
	}
	return nil
}

func listCredentials(local credentials.Manager, remote credentials.Manager) error {
	var creds []*credentials.Credentials = make([]*credentials.Credentials, 0)
	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		remoteCreds, err := remote.List()
		if err != nil {
			return err
		}
		creds = append(creds, remoteCreds...)
	}
	if flagLocal {
		localCreds, err := local.List()
		if err != nil {
			return err
		}
		creds = append(creds, localCreds...)
	}
	displayListCredentials(creds)
	return nil
}

func displayListCredentials(creds []*credentials.Credentials) {
	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tALIAS\tHOST\tPORT\tUSER\tSRC")
	for idx, cred := range creds {
		fmt.Fprintf(
			tw,
			"%d\t%s\t%s\t%d\t%s\t%s\n",
			idx,
			cred.Alias,
			cred.Host,
			cred.Port,
			cred.Username,
			cred.Source,
		)
	}
	tw.Flush()
}

func deleteCredentials(local credentials.Manager, remote credentials.Manager) error {
	if flagLocal {
		if err := local.Del(flagAlias); err != nil {
			return err
		}
		fmt.Println("local credentials deleted for", flagAlias)
	}
	if flagRemote {
		if remote == nil {
			return errRemoteNotConfigured
		}
		if err := remote.Del(flagAlias); err != nil {
			return err
		}
		fmt.Println("remote credentials deleted for", flagAlias)
	}
	return nil
}
