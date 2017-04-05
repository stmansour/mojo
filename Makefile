DIRS = util db scrapers admin mojosrv
RELDIR = ./tmp/mojo

mojo:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	go clean
	rm -rf mojo tmp

package:
	mkdir -p ./tmp/mojo
	for dir in $(DIRS); do make -C $$dir package;done

all: clean mojo package stats

stats:
	@echo "GO SOURCE CODE STATISTICS"
	@echo "----------------------------------------"
	@find . -name "*.go" | srcstats
	@echo "----------------------------------------"
