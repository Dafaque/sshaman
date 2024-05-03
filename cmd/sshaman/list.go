package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Dafaque/sshaman/internal/credentials"
)

func listCredentials(local credentials.Manager, remote credentials.Manager) error {
	creds, err := getCredentialsList(local, remote)
	if err != nil {
		return err
	}
	displayListCredentials(creds)
	return nil
}

func getCredentialsList(local credentials.Manager, remote credentials.Manager) ([]*credentials.Credentials, error) {
	var creds []*credentials.Credentials = make([]*credentials.Credentials, 0)
	if flagRemote {
		if remote == nil {
			return nil, errRemoteNotConfigured
		}
		remoteCreds, err := remote.List()
		if err != nil {
			return nil, err
		}
		creds = append(creds, remoteCreds...)
	}
	if flagLocal {
		localCreds, err := local.List()
		if err != nil {
			return nil, err
		}
		creds = append(creds, localCreds...)
	}
	return creds, nil
}

func displayListCredentials(creds []*credentials.Credentials) {
	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tALIAS\tHOST\tPORT\tUSER\tSOURCE")
	for idx, cred := range creds {
		var source string
		if cred.Source != nil {
			source = *cred.Source
		} else {
			source = "local"
		}
		fmt.Fprintf(
			tw,
			"%d\t%s\t%s\t%d\t%s\t%s\n",
			idx,
			cred.Alias,
			cred.Host,
			cred.Port,
			cred.Username,
			source,
		)
	}
	tw.Flush()
}
