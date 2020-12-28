# KUBEPAAS
### **KUBEPAAS let you Focus on application rather than managing kubernetes.**

Deploy your web app on kubernetes with ease, no more YAML to write by hand, take advantage of automation with well defined API.

![kubepaas-flow-diagram](https://github.com/urvil38/kubepaas-cli/blob/master/doc/images/kubepaas-flow.png)

# Usage

```

██╗  ██╗██╗   ██╗██████╗ ███████╗██████╗  █████╗  █████╗ ███████╗
██║ ██╔╝██║   ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝
█████╔╝ ██║   ██║██████╔╝█████╗  ██████╔╝███████║███████║███████╗
██╔═██╗ ██║   ██║██╔══██╗██╔══╝  ██╔═══╝ ██╔══██║██╔══██║╚════██║
██║  ██╗╚██████╔╝██████╔╝███████╗██║     ██║  ██║██║  ██║███████║
╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝

A tool for interacting with kubepaas platform
and used for all kind of command that This plateform will support

Usage:
  kubepaas [flags]
  kubepaas [command]

Available Commands:
  config      view and edit Kubepaas CLI properties
  deploy      Deploy the application to kubepaas platform
  help        Help about any command
  login       Log in to kubepaas platform
  logout      Log out from kubepaas platform
  logs        Print the logs of current running containers
  profile     A brief description of your command
  rollout     RollBack to Older or Newer Version
  signup      Sign up for kubepaas platform
  update      Used for changing values of diffrent configuration

Flags:
  -h, --help   help for kubepaas

Use "kubepaas [command] --help" for more information about a command.
```

# Download

- Download appropriate pre-compiled binary from the [release](https://github.com/urvil38/kubepaas/releases) page.

```
# download binary using cURL
$ curl -L https://github.com/urvil38/kubepaas/releases/download/0.0.1/kubepaas-darwin-amd64 -o kubepaas

# make binary executable
$ chmod +x ./kubepaas

# move it to bin dir (user need to has root privileges. run following command as root user using sudo.
$ sudo mv ./kubepaas /usr/local/bin
```


- Download using `go get`

```
$ go get -u github.com/urvil38/kubepaas
```

# Build

- If you want to build kubepaas right away, you need a working [Go environment](https://golang.org/doc/install). It requires Go version 1.12 and above.

```
$ git clone https://github.com/urvil38/kubepaas.git
$ cd kubepaas
$ make build
```

# DEMO

[![asciicast](https://asciinema.org/a/241272.svg)](https://asciinema.org/a/241272?t=8)
