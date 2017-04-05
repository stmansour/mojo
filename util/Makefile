TOP=..
THISDIR="util"

util: *.go
	go vet
	golint
	go build
	go install
	@echo "*** ${THISDIR}: completed util ***"

clean:
	go clean
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

package:
	@echo "*** util completed package ***"

all: clean util
