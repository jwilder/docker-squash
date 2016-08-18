docker-squash
=============

Squash docker images to make them smaller.

docker-squash is a utility to squash multiple docker layers into one in order to create an image
with fewer and smaller layers.  It retains Dockerfile commands such as PORT,
ENV, etc.. so that squashed images work the same as they were originally built.  In addition, deleted
files in later layers are actually purged from the image when squashed.

It's designed to support a workflow where you would squash the image just before pushing it
to a registry.  Before squashing the image, you would remove any build time dependencies, extra
files (apt caches, logs, private keys, etc..) that you would not want to deploy.  The defaults also
preserve your base image so that its contents are not repeatedly transferred when pushing and pulling
images.

Typical workflow is:

* Build image from Dockerfile (or other)
* Run clean up commands in container (apt-get purge, rm /path, etc..)
* Commit container
* Squash the image
* Push squashed image to registry

See [Squashing Docker Images](http://jasonwilder.com/blog/2014/08/19/squashing-docker-images/)

## Sample Output

```
$ docker save 49b5a7a88d5 | sudo docker-squash -t squash -verbose | docker load
Loading export from STDIN using /tmp/docker-squash683466637 for tempdir
Loaded image w/ 15 layers
Extracting layers...
  -  /tmp/docker-squash683466637/49b5a7a88d5353fe77204ad5591a3ef100fc2807a9d6dce979fd1b17a73a68d6/layer.tar
  -  /tmp/docker-squash683466637/651626d6e364ccc22ac990ba95cd0aab9256c56055087cc9a5a1790cea5250b9/layer.tar
  -  /tmp/docker-squash683466637/c50f2b65cab3b74f9bdb6f616b36f132b9a182ed883d03f11173e32fa39ab599/layer.tar
  -  /tmp/docker-squash683466637/d92c3c92fa73ba974eb409217bb86d8317b0727f42b73ef5a05153b729aaf96b/layer.tar
  -  /tmp/docker-squash683466637/511136ea3c5a64f264b78b5433614aec563103b4d4702f3ba7d4d2698e22c158/layer.tar
  -  /tmp/docker-squash683466637/1c9383292a8ff4c4196ff4ffa36e5ff24cb217606a8d1f471f4ad27c4690e290/layer.tar
  -  /tmp/docker-squash683466637/589338fba5eb5cc32a25a036975a5e0938f12eff0dc70b661363c13ef1a192a5/layer.tar
  -  /tmp/docker-squash683466637/63e174c2ca3d53e2b7639a440940e16e15c1970e6ad16f740ffdcc60e59e0324/layer.tar
  -  /tmp/docker-squash683466637/9942dd43ff211ba917d03637006a83934e847c003bef900e4808be8021dca7bd/layer.tar
  -  /tmp/docker-squash683466637/0ea0d582fd9027540c1f50c7f0149b237ed483d2b95ac8d107f9db5a912b4240/layer.tar
  -  /tmp/docker-squash683466637/8dfc0bb00563dab615dfcc28ab3e338089f5b1d71d82d731c18cbe9f7667435f/layer.tar
  -  /tmp/docker-squash683466637/c4ff7513909dedf4ddf3a450aea68cd817c42e698ebccf54755973576525c416/layer.tar
  -  /tmp/docker-squash683466637/cc58e55aa5a53b572f3b9009eb07e50989553b95a1545a27dcec830939892dba/layer.tar
  -  /tmp/docker-squash683466637/e4eea4411c0065f8b0c7cf6be31dd58daa5ac04d8c64d54537cbfce2eb8e3413/layer.tar
  -  /tmp/docker-squash683466637/fc294d2b22cb53cb2440ff6fece18813ee7363f5198f5e20346abfcf7cce54fe/layer.tar
Inserted new layer 27935276f797 after 1c9383292a8f
  -  511136ea3c5a
  -  1c9383292a8f /bin/sh -c #(nop) ADD file:c1472c26527df28498744f9e9e8a8304c
  -> 27935276f797 /bin/sh -c #(squash) from 1c9383292a8f
  -  9942dd43ff21 /bin/sh -c echo '#!/bin/sh' > /usr/sbin/policy-rc.d  && echo
  -  d92c3c92fa73 /bin/sh -c rm -rf /var/lib/apt/lists/*
  -  0ea0d582fd90 /bin/sh -c sed -i 's/^#\s*\(deb.*universe\)$/\1/g' /etc/apt/
  -  cc58e55aa5a5 /bin/sh -c apt-get update && apt-get dist-upgrade -y && rm -
  -  c4ff7513909d /bin/sh -c #(nop) CMD [/bin/bash]
  -  fc294d2b22cb /bin/sh -c apt-get update && apt-get install -y golang-go
  -  8dfc0bb00563 /bin/sh -c #(nop) ADD dir:78239d85b32dd28e4cb1d81ace7ffd32b8
  -  651626d6e364 /bin/sh -c #(nop) WORKDIR /app
  -  589338fba5eb /bin/sh -c go build -o http
  -  c50f2b65cab3 /bin/sh -c #(nop) ENV PORT=8000
  -  e4eea4411c00 /bin/sh -c #(nop) EXPOSE map[8000/tcp:{}]
  -  63e174c2ca3d /bin/sh -c #(nop) CMD [/app/http]
  -  49b5a7a88d53 /bin/bash
Squashing from 27935276f797 into 27935276f797
  -  Deleting whiteouts
  -  Rewriting child history
  -  Removing 9942dd43ff21. Squashed. (/bin/sh -c echo '#!/bin/sh' > /usr/sbin/policy-...)
  -  Removing d92c3c92fa73. Squashed. (/bin/sh -c rm -rf /var/lib/apt/lists/*)
  -  Removing 0ea0d582fd90. Squashed. (/bin/sh -c sed -i 's/^#\s*\(deb.*universe\)$/\1...)
  -  Removing cc58e55aa5a5. Squashed. (/bin/sh -c apt-get update && apt-get dist-upgra...)
  -  Replacing c4ff7513909d w/ new layer 72391e640b52 (/bin/sh -c #(nop) CMD [/bin/bash])
  -  Removing fc294d2b22cb. Squashed. (/bin/sh -c apt-get update && apt-get install -y...)
  -  Removing 8dfc0bb00563. Squashed. (/bin/sh -c #(nop) ADD dir:78239d85b32dd28e4cb1d...)
  -  Replacing 651626d6e364 w/ new layer bd7b4b11874a (/bin/sh -c #(nop) WORKDIR /app)
  -  Removing 589338fba5eb. Squashed. (/bin/sh -c go build -o http)
  -  Replacing c50f2b65cab3 w/ new layer e4af8871b961 (/bin/sh -c #(nop) ENV PORT=8000)
  -  Replacing e4eea4411c00 w/ new layer 6803497b6a61 (/bin/sh -c #(nop) EXPOSE map[8000/tcp:{}])
  -  Replacing 63e174c2ca3d w/ new layer 40b8c7c33bba (/bin/sh -c #(nop) CMD [/app/http])
  -  Removing 49b5a7a88d53. Squashed. (/bin/bash)
Tarring up squashed layer 27935276f797
Removing extracted layers
Tagging 40b8c7c33bba as jwilder/whoami:squash
Tarring new image to STDOUT
Done. New image created.
  -  40b8c7c33bba Less than a second /bin/sh -c #(nop) CMD [/app/http] 3.072 kB
  -  6803497b6a61 Less than a second /bin/sh -c #(nop) EXPOSE map[8000/tcp:{}] 3.072 kB
  -  e4af8871b961 Less than a second /bin/sh -c #(nop) ENV PORT=8000 3.072 kB
  -  bd7b4b11874a Less than a second /bin/sh -c #(nop) WORKDIR /app 3.072 kB
  -  72391e640b52 Less than a second /bin/sh -c #(nop) CMD [/bin/bash] 3.072 kB
  -  27935276f797 1 seconds /bin/sh -c #(squash) from 1c9383292a8f 39.49 MB
  -  1c9383292a8f 3 days /bin/sh -c #(nop) ADD file:c1472c26527df28498744f9e9e8a83... 201.6 MB
  -  511136ea3c5a 14 months  1.536 kB
Removing tempdir /tmp/docker-squash683466637
```

## Installation

Download the latest version:

* [Linux linux/amd64](https://github.com/jwilder/docker-squash/releases/download/v0.2.0/docker-squash-linux-amd64-v0.2.0.tar.gz)
* [OSX darwin/amd64](https://github.com/jwilder/docker-squash/releases/download/v0.2.0/docker-squash-darwin-amd64-v0.2.0.tar.gz)

```
$ wget https://github.com/jwilder/docker-squash/releases/download/v0.2.0/docker-squash-linux-amd64-v0.2.0.tar.gz
$ sudo tar -C /usr/local/bin -xzvf docker-squash-linux-amd64-v0.2.0.tar.gz
```
NOTE: docker-squash must run as root (to maintain file permission created within images).

Dependencies:

* [tar 1.27](http://www.gnu.org/software/tar/)

## Usage

docker-squash works by squashing a saved image and loading the squashed image back into docker.

```
$ docker save <image id> > image.tar
$ sudo docker-squash -i image.tar -o squashed.tar
$ cat squashed.tar | docker load
$ docker images <new image id>
```

You can also tag the squashed image:

```
$ docker save <image id> > image.tar
$ sudo docker-squash -i image.tar -o squashed.tar -t newtag
$ cat squashed.tar | docker load
$ docker images <new image id>
```

You can reduce disk IO by piping the input and output to and from docker:

```
$ docker save <image id> | sudo docker-squash -t newtag | docker load
```

If you have a sufficient amount of RAM, you can also use a `tmpfs` to remove temporary
disk storage:

```
$ docker save <image_id> | sudo TMPDIR=/var/run/shm docker-squash -t newtag | docker load
```

By default, a squashed layer is inserted after the first `FROM` layer.  You can specify a different
layer with the `-from` argument.
```
$ docker save <image_id> | sudo docker-squash -from <other layer> -t newtag | docker load
```
If you are creating a base image or only want one final squashed layer, you can use the
`-from root` to squash the base layer and your changes into one layer.

```
$ docker save <image_id> | sudo docker-squash -from root -t newtag | docker load
```

### Development

This project uses [glock](https://github.com/robfig/glock) for managing 3rd party dependencies.
You'll need to install glock into your workspace before hacking on docker-squash.

```
$ git clone <your fork>
$ glock sync github.com/jwilder/docker-squash
$ make
```

## License

MIT
