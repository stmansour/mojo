DIRS = ws mailsend msgbody

db:
	cd mailsend;make
	@echo "*** MAKE completed in test ***"

clean:
	rm -f testreport.txt
	for dir in $(DIRS); do make -C $$dir clean;done
	cd ./testdb;make clean
	cd ./mailsend;make clean
	@echo "*** CLEAN completed in test ***"

test:
	@echo "FUNCTIONAL TEST Report" > testreport.txt
	@echo "Report Start: \c" >> testreport.txt
	@date >> testreport.txt
	@echo >> testreport.txt
	@echo "Result  Path                  Test Name                                 Phases" >> testreport.txt
	@echo "------  --------------------  ----------------------------------------  ------" >> testreport.txt
	for dir in $(DIRS); do make -C $$dir test;done
	@./totals.sh
	@echo >> testreport.txt
	@echo "Report Complete: \c" >> testreport.txt
	@date >> testreport.txt
	@echo "*** TEST completed in test ***"
	@cat testreport.txt

package:
	for dir in $(DIRS); do make -C $$dir package;done
	cd mailsend;make package
	@echo "*** PACKAGE completed in test ***"
secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json
