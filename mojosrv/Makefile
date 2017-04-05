TOP=..
RELDIR=${TOP}/tmp/mojo
THISDIR="mojosrv"

mojosrv: *.go ver.go
	go vet
	golint
	go build
	@echo "*** ${THISDIR}: completed mojosrv ***"

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f mojosrv, ver.go, tmp *.log
	@echo "*** ${THISDIR}: completed clean ***"

package:
	mkdir -p ${RELDIR}
	cp mojosrv ${RELDIR}/mojosrv
	@echo "*** ${THISDIR}: completed package ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mojosrv package
