DIRS=mojonewdb

admin:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done

install:
	for dir in $(DIRS); do make -C $$dir install;done

package:
	for dir in $(DIRS); do make -C $$dir package;done

all: clean admin
