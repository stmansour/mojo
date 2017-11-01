#!/bin/bash
rm -rf ./csvdump
mkdir ./csvdump

mysqldump --no-defaults --tab=./csvdump --fields-terminated-by=, --fields-enclosed-by='"' --lines-terminated-by=0x0d0a mojo

cd ./csvdump
cat >p <<EOF
PID,FirstName,MiddleName,LastName,PreferredName,JobTitle,OfficePhone,OfficeFax,Email1,Email2,MailAddress,MailAddress2,MailCity,MailState,MailPostalCode,MailCountry,RoomNumber,MailStop,Status,OptOutDate,LastModTime,LastModBy
EOF

cat p People.txt > People.csv
open People.csv
