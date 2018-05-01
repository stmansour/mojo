DIRS = charnorm
THIS=tools

tools:
	for dir in $(DIRS); do make -C $$dir;done
	@echo "*** MAKE completed in ${THIS} ***"

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	@echo "*** CLEAN completed in ${THIS} ***"

test:
	@echo "*** TEST completed in ${THIS} ***"

package:
	for dir in $(DIRS); do make -C $$dir package;done
	@echo "*** PACKAGE completed in ${THIS} ***"
secure:
	for dir in $(DIRS); do make -C $${dir} secure;done
	@rm -f config.json confdev.json confprod.json
