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

try: clean mojo package

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
	@echo "recreating databases used in testing..."
	@echo "Getting people with last name Aa* from FAA"
	cd scrapers/faa;make q
	@echo "Adding sendmail test info"
	cd test/sendmail;./sendmail -n
	@echo "Setting updated small testdb to db used in ./test/ws"
	cd test/testdb;make snapshot;make test
	@echo "Completed."

smalldb:
	@echo "making smalldb..."
	cd test/testdb;make smalldb
