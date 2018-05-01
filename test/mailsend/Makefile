TOP=../..
BINDIR=${TOP}/tmp/mojo
THISDIR="mailsend"
COUNTOL=${TOP}/tools/bashtools/countol.sh

mailsend: *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	@${COUNTOL} "go vet"
	@${COUNTOL} golint
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

config.json:
	/usr/local/accord/bin/getfile.sh accord/db/confdev.json
	cp confdev.json config.json

clean:
	go clean
	rm -f mailsend, ver.go, tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

package: mailsend
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	cp mailsend ${BINDIR}/mailsend
	cp testmsg.html mscript.sh resetdb.sh ${BINDIR}/
	@rm -f fail
	@echo "*** ${THISDIR}: completed package ***"

pubcamp:
	tar cvfz campaigns.tar campaigns
	scp -i ~/.ssh/smanAWS1.pem campaigns.tar dir3:~/apps/mojo
	ssh ec2-user@dir3 'cd apps/mojo;tar xzvf campaigns.tar'
	scp -i ~/.ssh/smanAWS1.pem campaigns.tar dir4:~/apps/mojo
	ssh ec2-user@dir4 'cd apps/mojo;tar xzvf campaigns.tar'

newdb:
	${BINDIR}/mojonewdb
	./mailsend -setup

tar:
	tar cvf campaigns.tar campaigns

untar:
	tar xvf campaigns.tar

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean mailsend package
secure:
	@rm -f config.json confdev.json confprod.json
