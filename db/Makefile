TOP=..
DIRS=schema
THISDIR=db

db: *.go
	for dir in $(DIRS); do make -C $$dir;done
	go vet
	golint
	go build
	go install
	@echo "*** ${THISDIR}: completed db ***"

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	go clean
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"


package:
	go install
	for dir in $(DIRS); do make -C $$dir package;done
	@echo "*** ${THISDIR}: completed package ***"

all: clean db package
