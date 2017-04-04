DIRS = util db scrapers admin mojosrv
RELDIR = ./tmp/mojo

mojo:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	go clean
	rm -f mojo

install:
	for dir in $(DIRS); do make -C $$dir install;done

all: clean mojo
