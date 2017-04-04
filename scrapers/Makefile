DIRS=faa

scrapers:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done

install:

all: clean scrapers
