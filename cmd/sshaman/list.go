package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Dafaque/sshaman/internal/credentials"
)

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
