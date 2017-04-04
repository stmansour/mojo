DIRS=mojofaa

admin:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done

install:
	for dir in $(DIRS); do make -C $$dir install;done

all: clean admin
