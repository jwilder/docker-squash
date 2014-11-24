First make sure that go is installed:

```bash
brew install go
mkdir $HOME/.go
export GOPATH=$HOME/.go
export PATH=$PATH:$GOPATH/bin
```

Clone (not download) the repository:

```bash
git clone https://github.com/jwilder/docker-squash.git
```

Install the units package:

```bash
go get github.com/docker/docker/pkg/units
```

Build the binary:

```bash
GOBIN=$(pwd) make
```

Now you should have a `docker-squash` executable in the current folder. Only one step missing to run it correctly:

Install GNU tar (you could also patch utils.go(in `extractTar` the `--xattrs` attribute is not supported in osx) but this seems easier to do)

```bash
brew install gnu-tar
PATH="/usr/local/opt/gnu-tar/libexec/gnubin:$PATH"
```

You can copy the `docker-squash` executable to whatever folder you are working on or put in in some place in your path (one of the many `bin` folders)
If you close the terminal and want to use it again in the future don't forget to run again all the exports commands, so that the environment variabels are 
correctly set. If you are using the same executable (not building again) then the only important export would be the last one (for GNU tar).
