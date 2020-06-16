# Emacs Cleaner `0.1.0`

Remove emacs temporary files

## Usage

Simple utility.
``` shell
$ emacs-cleaner .
$ emacs-cleaner --threads 20 .
$ emacs-cleaner --dry .
```

Default amount of threads is `10`. If given `--dry` will only print results and won't unlink anything. 
Give `--threads` the value `1` for single threaded usage (recommended for slow disks, strange filesystems, or NFS mounts.)

## Building

To just run:
``` shell
$ go run emacs-cleaner.go
```

To just build the binary:
``` shell
$ go build emacs-cleaner.go
```

To build and install:
``` shell
$ make && sudo make install
```
The default install path is `/usr/local/bin`, change `INSTALL_DIR` in [Makefile] to specify location.

 [Makefile]: ./blob/master/Makefile
 
To uninstall:
``` shell
$ sudo make uninstall
```

## Liscense
GPL'd with love <3
