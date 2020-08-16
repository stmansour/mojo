#!/bin/bash
ssh ec2-user@dir3 'rm -f mojo.sql.gz;/usr/bin/mysqldump -h phbk.cjkdwqbdvxyu.us-east-1.rds.amazonaws.com -P 3306 mojo > mojo.sql; gzip mojo.sql'
scp -i ~/.ssh/smanAWS1.pem dir3:~/mojo.sql.gz .
gunzip -f mojo.sql.gz
