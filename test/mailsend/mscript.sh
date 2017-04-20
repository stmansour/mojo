#!/bin/bash

#-----------------------------------------------
#  Development:  Steve  +  Amazon Test Accounts
#-----------------------------------------------
#./mailsend -h "http://localhost:8275/" -q MojoTest -subject "Big Mojo Test" -a perks.pdf
#./mailsend -h "http://localhost:8275/" -q MojoTest -subject "Perks Mojo Test" -a perks.pdf -b msg.html

./mailtester -h "http://ec2-54-152-108-202.compute-1.amazonaws.com:8275/" -q MojoTest -subject "Big Mojo Test" -a perks.pdf

#./mailtester -h "http://ec2-54-152-108-202.compute-1.amazonaws.com:8275/" -q MojoTest -subject "Perks - Mojo Email Test" -a perks.pdf -b msg.html



#--------------------------------------------
#  AMAZON:  Steve  +  Amazon Test Accounts
#--------------------------------------------
#./mailtester -h "http://ec2-54-152-108-202.compute-1.amazonaws.com:8275/" -q AmazonTest -subject "Big Mojo Test 3" -a perks.pdf


#------------------------------------------------------
#  ACCORD:  Steve, Joe, Melissa,  Amazon Test Accounts
#------------------------------------------------------
#./mailtester -h "http://ec2-54-152-108-202.compute-1.amazonaws.com:8275/" -q AccordTest -subject "Big Mojo Test to Accord" -a perks.pdf
