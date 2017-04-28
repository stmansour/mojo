#!/bin/bash

#-----------------------------------------------
#  Development:  Steve  +  Amazon Test Accounts
#-----------------------------------------------
#./mailsend -b ftmsg.html -subject "FAA Today"

#--------------------------------------------
#  AMAZON:  Steve  +  Amazon Test Accounts
#--------------------------------------------
#./mailsend -b ftmsg.html -subject "FAA Today" -q AmazonTest

#------------------------------------------------------
#  ACCORD:  Steve, Joe, Melissa,  Amazon Test Accounts
#------------------------------------------------------
#./mailsend -b ftmsg.html -subject "FAA Today" -q AccordTest

#------------------------------------------------------
#  FAA: All of FAA
#------------------------------------------------------
./mailsend -from "Editor-in-chief@FAAToday.com" -b ftmsg.html -subject "FAA Today"



# ./mailsend -from "Editor-in-chief@FAAToday.com" -b ftmsg.html -subject "FAA Today" -q "FAA-1-First50"
# ./mailsend -from "Editor-in-chief@FAAToday.com" -b ftmsg.html -subject "FAA Today" -q "FAA-2-Next250"
# ./mailsend -from "Editor-in-chief@FAAToday.com" -b ftmsg.html -subject "FAA Today" -q "FAA-3-Next700"
# ./mailsend -from "Editor-in-chief@FAAToday.com" -b ftmsg.html -subject "FAA Today" -q "FAA-4-Next5000"
# ./mailsend -from "Editor-in-chief@FAAToday.com" -b ftmsg.html -subject "FAA Today" -q "FAA-5-Next20000"
# ./mailsend -from "Editor-in-chief@FAAToday.com" -b ftmsg.html -subject "FAA Today" -q "FAA-6-TheRest"
