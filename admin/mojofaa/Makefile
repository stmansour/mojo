TOP=../..
RELDIR=${TOP}/tmp/mojo

mojofaa: *.go
	go vet
	golint
	go build

clean:
	go clean
	rm -f mojofaa

install:
	mkdir -p ${RELDIR}
	cp mojofaa ${RELDIR}/mojofaa


all: clean mojofaa
