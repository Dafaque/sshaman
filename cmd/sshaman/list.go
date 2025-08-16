package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Dafaque/sshaman/v2/internal/credentials"
)

func listCredentials(manager *credentials.Manager) error {
	creds, err := getCredentialsList(manager)
	if err != nil {
		return err
	}
	displayListCredentials(creds)
	return nil
}

func getCredentialsList(manager *credentials.Manager) ([]*credentials.Credentials, error) {
	creds, err := manager.List()
	if err != nil {
		return nil, err
	}
	return creds, nil
}

func displayListCredentials(creds []*credentials.Credentials) {
	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tNAME\tHOST\tPORT\tUSER")
	for idx, cred := range creds {
		fmt.Fprintf(
			tw,
			"%d\t%s\t%s\t%d\t%s\n",
			idx,
			cred.Name,
			cred.Host,
			cred.Port,
			cred.UserName,
		)
	}
	tw.Flush()
}
