# SSHaMan
SSHaMan is open-source SSH connections manager written in go.

# Install
```shell
go install github.com/Dafaque/sshaman/cmd/sshaman@latest
```
# Config
## Environment variables
- `SSHAMAN_HOME` - path to the directory where the application will store its data
# Usage
## Add credentials
Command bellow will help you create new credentials for `localhost:22` with user `admin` and ssh key `~/.ssh/id_ed25519`:
```shell
sshaman add \
    -n myfirstserver \
    -h localhost \
    -u admin \
    -k ~/.ssh/id_ed25519 \
    -no-pw \
    -no-pp
```
## Edit credentils
Commnd bellow will override `myfirstserver` with user `user` and password authentication:
```shell
sshaman add \
    -n myfirstserver \
    -h localhost \
    -u user \
    -f
```
## List credentials
```shell
sshaman list
```
## Connect using created alias
```shell
sshaman connect -n myfirstserver
```
## Export credentials
To create `sshaman.enc` in the current directory run:
```shell
sshaman export
```
> [!WARNING]
> You can skip the password prompt, but this is not recommended since the application is open source

> [!CAUTION]
> Exporting credentials will overwrite the existing `sshaman.enc` file.

## Import credentials
To import recently exported `sshaman.enc` file `cd` to the directory where the file is located and run:
```shell
sshaman import
```
## Delete credentials
Command below will delete credentials
```shell
sshaman remove -n myfirstserver
```
> [!CAUTION]
> This operation is irreversible.
## Delete all credentials
Command below will delete credentials
```shell
sshaman drop
```
> [!CAUTION]
> This operation is irreversible.

## All commands
```
Usage of add:
  -f    force operation
  -h string
        new ssh connection's address
  -k string
        new ssh connection's key file path
  -n string
        new ssh connection's alias
  -no-pp
        skip key's passphrase prompt
  -no-pw
        skip password prompt
  -p int
        new ssh connection's port (default 22)
  -u string
        new ssh connection's user
Usage of connect:
  -n string
        ssh connection's allias to conect
Usage of remove:
  -n string
        ssh connection's allias to conect
Usage of drop:
  -f    force operation
Usage of import:
  -dry
        view what would be imported
  -no-pw
        skip password prompt
Usage of export:
  -no-pw
        skip password prompt
Usage sshaman:
  -v    show app details
```
