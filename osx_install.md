First make sure that go is installed:

```bash
brew install go
mkdir $HOME/.go
export GOPATH=$HOME/.go
export PATH=$PATH:$GOPATH/bin
```

Clone this repository into your GOPATH src folder

```bash
mkdir -p $GOPATH/src/github.com/jwilder
cd $GOPATH/src/github.com/jwilder
git clone https://github.com/jwilder/docker-squash.git
```

Install [glock](https://github.com/robfig/glock) as docker-squash uses glock for managing dependencies. This project depends on [docker/docker@2606a2](https://github.com/docker/docker/tree/2606a2e4d3bf810ec82e373a6cd334e22e504e83) which contains [pkg/units](https://github.com/docker/docker/tree/2606a2e4d3bf810ec82e373a6cd334e22e504e83/pkg/units).

```bash
go get github.com/robfig/glock
```

Sync dependencies and install units package

```bash
glock sync github.com/jwilder/docker-squash
go get github.com/docker/docker/pkg/units
```

In the docker-squash repo, build the binary:

```bash
cd $GOPATH/src/github.com/jwilder/docker-squash
make
```

Now you should have a `docker-squash` executable in $GOPATH/bin. Only one step missing to run it correctly:

Install GNU tar (you could also patch utils.go(in `extractTar` the `--xattrs` attribute is not supported in osx) but this seems easier to do)

```bash
brew install gnu-tar
PATH="/usr/local/opt/gnu-tar/libexec/gnubin:$PATH"
```

You can copy the `docker-squash` executable to whatever folder you are working on or put in in some place in your path (one of the many `bin` folders)
If you close the terminal and want to use it again in the future don't forget to run again all the exports commands, so that the environment variabels are
correctly set. If you are using the same executable (not building again) then the only important export would be the last one (for GNU tar).
