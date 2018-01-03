#!/bin/bash

# Fix common email address errors

declare -a out_filters=(
	'@gmail$'
	'([^@])gmail.com'
	'([^@])msn.com'
	'([^@])yahoo.com'
	'@verizon$'
	'@gmail$'
	'@yahoo$'
	'@msn$'
	"'([a-zA-Z.]+)@faa.gov"
	'@gmailcom$'
	'@yahoocom$'
	'@msncom$'
	'@verizoncom$'
	'@faagov$'
	'@faa$'
	'@faa,gov$'
	'@gmail,com$'
	'@yahoo,com$'
	'@msn,com$'
	'@verizon,com'
)
