#!/bin/bash

#---------------------------------------------------------------
# TOP is the directory where RentRoll begins. It is used
# in base.sh to set other useful directories such as ${BASHDIR}
#---------------------------------------------------------------
TOP=../..

TESTNAME="Web Services"
TESTSUMMARY="Test Web Services"

CREATENEWDB=0

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
# echo "Create new database..." 
# mysql --no-defaults rentroll < restore.sql

source ../share/base.sh

echo "STARTING MOJO SERVER"
startMojoServer

curl http://localhost:8275/v1/ping/


# Get Specificy PaymentType 
# echo "request%3D%7B%22cmd%22%3A%22get%22%2C%22selected%22%3A%5B%5D%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
# dojsonPOST "http://localhost:8275/v1/people/1" "request" "y"  "WebService--GetPeople"
# dojsonGET "http://localhost:8275/v1/peoplecount/" "a"  "WebService--PeopleCount"
# dojsonGET "http://localhost:8275/v1/groupcount/" "b"  "WebService--GroupCount"
# dojsonGET "http://localhost:8275/v1/peoplestats/" "c"  "WebService--PeopleStats"

# echo "request%3D%7B%22cmd%22%3A%22%22%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
# dojsonPOST "http://localhost:8275/v1/queries/" "request" "d"  "WebService--SearchQueries"

# dojsonGET "http://localhost:8275/v1/qrescount/1" "e"  "WebService--QueryResultsCount"

# # Subscription Succeeded Notification
# echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%22b7001127-9d1f-524a-9045-404e3913897b%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%22%2C%22Message%22%3A%7B%22notificationType%22%3A%22AmazonSnsSubscriptionSucceeded%22%2C%22message%22%3A%22You%20have%20successfully%20subscribed%20your%20Amazon%20SNS%20topic%20'arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email'%20to%20receive%20'Bounce'%20notifications%20from%20Amazon%20SES%20for%20identity%20'sman%40accordinterests.com'.%22%7D%2C%22Timestamp%22%3A%222017-04-13T16%3A05%3A32.839Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22Dv6o%20yiGfg3dzDkRTZ2Wf4Uj08DtzhziuPW99QkGDes9pLd3lOIUwl%20u5XEfURI374NiZD3pr7Ku7U7WwVCA%20Laa4RlGF8T5FXlpaLuDYgccEgoDNm1IAX%204a19yUdWD7dAzv1eEIcQSTgezT4QGYZh1XuTAvUAU%2FZnBcfL5P8nGAAKrf78jEUkV0rcH%2FHRw1ndewlzqcQoDy2H64x%20EqHuYIO%2F7mbTgDFQNxaQg0Fr64ipzheVfijvWe4NpWjblAAawuJRnPyEpDfFdJZFbfk%206y8QpC8O56QDVXV5h5PaKvioqo0KyOyK9HpuTGRgy%20X5OZAxjKtLXxwVZF%20wzdQ%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%3A4ed1756c-e31f-4809-a873-056c3026afb4%22%7D" > request
# dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "f"  "WebService--AWS-SubscriptionSucceeded" "Notification"

# Bounce Notification
echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%2239a64847-8102-5692-bcd7-192d924f075d%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%22%2C%22Message%22%3A%22%7B%5C%22notificationType%5C%22%3A%5C%22Bounce%5C%22%2C%5C%22bounce%5C%22%3A%7B%5C%22bounceType%5C%22%3A%5C%22Permanent%5C%22%2C%5C%22bounceSubType%5C%22%3A%5C%22General%5C%22%2C%5C%22bouncedRecipients%5C%22%3A%5B%7B%5C%22emailAddress%5C%22%3A%5C%22bounce%40simulator.amazonses.com%5C%22%2C%5C%22action%5C%22%3A%5C%22failed%5C%22%2C%5C%22status%5C%22%3A%5C%225.1.1%5C%22%2C%5C%22diagnosticCode%5C%22%3A%5C%22smtp%3B%20550%205.1.1%20user%20unknown%5C%22%7D%5D%2C%5C%22timestamp%5C%22%3A%5C%222017-04-14T00%3A00%3A13.068Z%5C%22%2C%5C%22feedbackId%5C%22%3A%5C%220100015b69c29ab8-fbf8d711-ea80-4465-bb85-aa901ea7f145-000000%5C%22%2C%5C%22remoteMtaIp%5C%22%3A%5C%22207.171.163.188%5C%22%2C%5C%22reportingMTA%5C%22%3A%5C%22dsn%3B%20a8-50.smtp-out.amazonses.com%5C%22%7D%2C%5C%22mail%5C%22%3A%7B%5C%22timestamp%5C%22%3A%5C%222017-04-14T00%3A00%3A08.000Z%5C%22%2C%5C%22source%5C%22%3A%5C%22sman%40stevemansour.com%5C%22%2C%5C%22sourceArn%5C%22%3A%5C%22arn%3Aaws%3Ases%3Aus-east-1%3A553053000801%3Aidentity%2Fsman%40stevemansour.com%5C%22%2C%5C%22sourceIp%5C%22%3A%5C%2224.6.191.18%5C%22%2C%5C%22sendingAccountId%5C%22%3A%5C%22553053000801%5C%22%2C%5C%22messageId%5C%22%3A%5C%220100015b69c28a6f-956eba47-89be-4ea8-8e5e-9924ee3954eb-000000%5C%22%2C%5C%22destination%5C%22%3A%5B%5C%22bounce%40simulator.amazonses.com%5C%22%5D%7D%7D%22%2C%22Timestamp%22%3A%222017-04-14T00%3A00%3A13.118Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22Akcg6IxxVQ1aIC0YIJMnuY9FVaNpndIQjuzVVFWVHZ7n1jKjIZYeH%2Bt1AAQBh%2FXcucIVlhKKfKjqgLTimPpde%2B4FEc6JS6G3eQW1pxkKuyDlPdj0AjUbMouCKNKldmFGWlA%2BuAF2GzXHUQPAwPMpnl%2FcyIySJgqlQk52ODfkzdJ1EFQM6D2jKPXrs1B5OdSgb9e3uGCo4R3U3wQkFbEVtj7%2BTbsv9AAEMwrVyrhhUDlBitFMpBVkI0Cu%2BVGf6V%2BefQtBe40l8t%2BEHVzdfDy2RYDH57PgnOVqENRR3NGsZ2ugqdKlLfJaVEgfhug6fi%2F4i%2FmQndXNhGLlgLzz8Fzl7Q%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%3A41abeec0-7f87-419f-a33e-fd97eda671bd%22%7D" > request
dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "g"  "WebService--AWS-BounceNotification" "Notification"


stopMojoServer
echo "RENTROLL SERVER STOPPED" 

logcheck
