package ssh

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"github.com/Dafaque/sshaman/internal/credentials"
)

type SshClient struct {
	Config *ssh.ClientConfig
	// @todo for future errors verbosing
	Server string
}

func NewSshClient(creds *credentials.Credentials) (*SshClient, error) {
	amc := NewAuthMethodConfig()
	if creds.Password != nil {
		amc = amc.WithPassword(*creds.Password)
	}
	if creds.Key != nil {
		amc = amc.WithKeyPassphrase(creds.Key, creds.Passphrase)
	}
	// build SSH client config
	auth, err := amc.AuthMethods()
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: creds.Username,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// @todo use OpenSSH's known_hosts file
			return nil
		},
	}
	client := &SshClient{
		Config: config,
		Server: fmt.Sprintf("%v:%v", creds.Host, creds.Port),
	}

	return client, nil
}

func (s *SshClient) Loop() error {
	conn, err := ssh.Dial("tcp", s.Server, s.Config)
	if err != nil {
		return fmt.Errorf("Dial to %v failed %v", s.Server, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("Create session for %v failed %v", s.Server, err)
	}
	defer session.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			var sessig ssh.Signal
			switch sig {
			case os.Interrupt:
				sessig = ssh.SIGINT
			default:
				break
			}
			if err := session.Signal(sessig); err != nil {
				// @todo chan
				panic(err)
			}
		}
	}()

	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("terminal make raw: %s", err)
	}
	defer term.Restore(fd, state)

	w, h, err := term.GetSize(fd)
	if err != nil {
		return fmt.Errorf("terminal get size: %s", err)
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("linux", h, w, modes); err != nil {
		return err
	}

	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		return err
	}
	return session.Wait()
}
