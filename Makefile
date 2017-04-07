DIRS = util db scrapers admin mojosrv test
RELDIR = ./tmp/mojo

.PHONY:  test

mojo:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	go clean
	rm -rf mojo tmp

package:
	mkdir -p ./tmp/mojo
	for dir in $(DIRS); do make -C $$dir package;done

test: package
	for dir in $(DIRS); do make -C $$dir test;done
	cat test/testreport.txt

all: clean mojo test stats

stats:
	@echo "GO SOURCE CODE STATISTICS"
	@echo "----------------------------------------"
	@find . -name "*.go" | srcstats
	@echo "----------------------------------------"
