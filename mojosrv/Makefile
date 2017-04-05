TOP=..
RELDIR=${TOP}/tmp/mojo

mojosrv: *.go ver.go
	go vet
	golint
	go build

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f mojosrv, ver.go, tmp *.log

package:
	mkdir -p ${RELDIR}
	cp mojosrv ${RELDIR}/mojosrv


all: clean mojosrv package
