SRC:= emacs-cleaner.go
INSTALL_DIR:= /usr/local/bin

all: clean emacs-cleaner

clean:
	rm -f emacs-cleaner

emacs-cleaner:
	go build $(SRC)
	strip $@

install:
	cp -f emacs-cleaner $(INSTALL_DIR)/emacs-cleaner

uninstall:
	rm -f $(INSTALL_DIR)/emacs-cleaner
