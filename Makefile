
all: clean emacs-cleaner

clean:
	rm -f emacs-cleaner

emacs-cleaner:
	go build emacs-cleaner.go
	strip $@
