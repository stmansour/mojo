DIRS=faa

scrapers:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done

package:
	for dir in $(DIRS); do make -C $$dir package;done

all: clean scrapers package
