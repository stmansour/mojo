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
	cp activate.sh ./tmp/mojo/
	cp update.sh ./tmp/mojo/
	for dir in $(DIRS); do make -C $$dir package;done

test: package
	for dir in $(DIRS); do make -C $$dir test;done

all: clean mojo test stats

try: clean mojo package smalldb

stats:
	@echo "GO SOURCE CODE STATISTICS"
	@echo "----------------------------------------"
	@find . -name "*.go" | srcstats
	@echo "----------------------------------------"

# Sometimes the database schema changes. When this happens many
# things won't work because to speed up testing we use mysql to 
# restore test databases -- and the restored databases will not
# have the correct schema to match the updated schema. So, use
# this target to regenerate the databases. The way to use this
# target is typically as follows:
# 	a) make try
#	b) make schemachange
#	c) make test
schemachange:
	@echo "Rebuilding test database..."
	cd scrapers/faa;make q
	cd test/testdb;make smalldb;make test
	@echo "Completed."

smalldb:
	@echo "making smalldb..."
	cd test/testdb;make smalldb

testdb: smalldb

publish: package
	cd tmp;tar cvf mojo.tar mojo; gzip mojo.tar
	cd tmp;/usr/local/accord/bin/deployfile.sh mojo.tar.gz jenkins-snapshot/mojo/latest
