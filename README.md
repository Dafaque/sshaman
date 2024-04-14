# SSHaMan
![sshaman logo by kandinsky 3.0](assets/logo_kandinsky3.0_256.jpg "SSHaMan")

SSHaMan is open-source SSH connections manager written in go.

# Install 

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
    -skip-passphrase \
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
## Delete credentials
Command below will delete credentials
```shell
sshaman delete -alias myfirstserver
```
## Delete all credentials
Command below will delete credentials
```shell
sshaman drop
```

## All commands
```shell
sshaman -h
```
```
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
Usage of sshaman:
  -alsologtostderr
        log to standard error as well as files
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -v value
        log level for V logs
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
```
# ToDo
- Add app/OpenSSH info
- OpenSSH's known_hosts handle
- App configuration
- Import/export repository
- Local credentials repository encryption
- Remote credentials repository