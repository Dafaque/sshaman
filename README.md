# SSHaMan
![sshaman logo by kandinsky 3.0](assets/logo_kandinsky3.0_256.jpg "SSHaMan")

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
    -alias myfirstserver \
    -host localhost \
    -user admin \
    -key ~/.ssh/id_ed25519 \
    -skip-password \
    -skip-passphrase
```
## Edit credentils
Commnd bellow will override `myfirstserver` with user `user` and password authentication:
```shell
sshaman add \
    -lias myfirstserver \
    -host localhost \
    -user user \
    -force
```
## List credentials
```shell
sshaman list
```
## Connect using created alias
```shell
sshaman connect -alias myfirstserver
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
sshaman delete -alias myfirstserver
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
```shell
sshaman -h
Usage of add:
  -alias string
        new ssh connection's alias
  -force
        force operation
  -host string
        new ssh connection's address
  -key string
        new ssh connection's key file path
  -local
        use local storage (default true)
  -port int
        new ssh connection's port (default 22)
  -remote
        use remote storage (unimplementer)
  -skip-passphrase
        skip key's passphrase prompt
  -skip-password
        skip password prompt
  -user string
        new ssh connection's user
Usage of connect:
  -alias string
        ssh connection's allias to conect
  -local
        use local storage (default true)
  -remote
        use remote storage (unimplemented)
Usage of list:
  -local
        use local storage (default true)
  -remote
        use remote storage (unimplemented)
Usage of delete:
  -alias string
        ssh connection's allias to conect
  -local
        use local storage (default true)
  -remote
        use remote storage (unimplemented)
Usage of drop:
  -force
        force operation
  -local
        use local storage (default true)
  -remote
        use remote storage (unimplemented)
Usage of import:
  -dry-run
        view what would be imported
  -skip-password
        skip password prompt
Usage of export:
  -skip-password
        skip password prompt
Usage of sshaman:
  -version
        show app details
```

# ToDo
- Add Homebrew release
- Add apt/snap release
- Add nix derivation
- a.k.a OpenSSH's known_hosts handle
- Remote credentials repository