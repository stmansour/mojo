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
echo "request%3D%7B%22cmd%22%3A%22get%22%2C%22selected%22%3A%5B%5D%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
dojsonPOST "http://localhost:8275/v1/people/1" "request" "y"  "WebService--GetPeople"
dojsonGET "http://localhost:8275/v1/peoplecount/" "a"  "WebService--PeopleCount"
dojsonGET "http://localhost:8275/v1/groupcount/" "b"  "WebService--GroupCount"
dojsonGET "http://localhost:8275/v1/peoplestats/" "c"  "WebService--PeopleStats"

echo "request%3D%7B%22cmd%22%3A%22%22%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
dojsonPOST "http://localhost:8275/v1/queries/" "request" "d"  "WebService--SearchQueries"

dojsonGET "http://localhost:8275/v1/qrescount/1" "e"  "WebService--QueryResultsCount"

# Subscription Succeeded Notification
echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%22b7001127-9d1f-524a-9045-404e3913897b%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%22%2C%22Message%22%3A%7B%22notificationType%22%3A%22AmazonSnsSubscriptionSucceeded%22%2C%22message%22%3A%22You%20have%20successfully%20subscribed%20your%20Amazon%20SNS%20topic%20'arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email'%20to%20receive%20'Bounce'%20notifications%20from%20Amazon%20SES%20for%20identity%20'sman%40accordinterests.com'.%22%7D%2C%22Timestamp%22%3A%222017-04-13T16%3A05%3A32.839Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22Dv6o%20yiGfg3dzDkRTZ2Wf4Uj08DtzhziuPW99QkGDes9pLd3lOIUwl%20u5XEfURI374NiZD3pr7Ku7U7WwVCA%20Laa4RlGF8T5FXlpaLuDYgccEgoDNm1IAX%204a19yUdWD7dAzv1eEIcQSTgezT4QGYZh1XuTAvUAU%2FZnBcfL5P8nGAAKrf78jEUkV0rcH%2FHRw1ndewlzqcQoDy2H64x%20EqHuYIO%2F7mbTgDFQNxaQg0Fr64ipzheVfijvWe4NpWjblAAawuJRnPyEpDfFdJZFbfk%206y8QpC8O56QDVXV5h5PaKvioqo0KyOyK9HpuTGRgy%20X5OZAxjKtLXxwVZF%20wzdQ%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%3A4ed1756c-e31f-4809-a873-056c3026afb4%22%7D" > request
dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "f"  "WebService--AWS-SubscriptionSucceeded" "Notification"

# Bounce Notification
echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%22c19a8ac9-8eb2-5ab3-aba2-dba024290846%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%22%2C%22Message%22%3A%7B%22notificationType%22%3A%22Bounce%22%2C%22bounce%22%3A%7B%22bounceType%22%3A%22Permanent%22%2C%22bounceSubType%22%3A%22General%22%2C%22bouncedRecipients%22%3A%5B%7B%22emailAddress%22%3A%22bounce%40simulator.amazonses.com%22%2C%22action%22%3A%22failed%22%2C%22status%22%3A%225.1.1%22%2C%22diagnosticCode%22%3A%22smtp%3B%20550%205.1.1%20user%20unknown%22%7D%5D%2C%22timestamp%22%3A%222017-04-13T16%3A33%3A31.929Z%22%2C%22feedbackId%22%3A%220100015b6829a6d7-22f25578-370e-4384-a927-bd75fff97c79-000000%22%2C%22remoteMtaIp%22%3A%22205.251.242.49%22%2C%22reportingMTA%22%3A%22dsn%3B%20a8-54.smtp-out.amazonses.com%22%7D%2C%22mail%22%3A%7B%22timestamp%22%3A%222017-04-13T16%3A33%3A31.000Z%22%2C%22source%22%3A%22sman%40stevemansour.com%22%2C%22sourceArn%22%3A%22arn%3Aaws%3Ases%3Aus-east-1%3A553053000801%3Aidentity%2Fsman%40stevemansour.com%22%2C%22sourceIp%22%3A%2224.6.191.18%22%2C%22sendingAccountId%22%3A%22553053000801%22%2C%22messageId%22%3A%220100015b6829a4ea-5b72ef43-19bb-41a7-8236-94286648b25b-000000%22%2C%22destination%22%3A%5B%22bounce%40simulator.amazonses.com%22%5D%7D%7D%2C%22Timestamp%22%3A%222017-04-13T16%3A33%3A31.962Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22ALo2Yn4c8I%20ynvYyyFGRTtOdkVMK4o3SIHKEcLg2p7OHT1qOrbrTf3YhTFZGIG2xqEZUufS7CjpULm3VWuwrDyjuZAGL0cg8%2FKqb%20NVNss1GtXYHuVNIuW5uLAmKA9SJgVVUXxjYQz%2FDuiQnslxet%20VdrFv0uY9vY21SShkkdi%20tQkcKSY3iiZeqF9%20eA40%20%20xrJXzrZtwSXT%20p0dFm8k4Rt4sxvMyi9meHqAB5cKaQE%20%2FPb7bY8f4gqKzbiuleLdEOubRwyIVdxsBq%2FdSUyKs8wlkjgjfJjmdEUvNWzoasazqaTsjh77ArXO3i9iRRdvdvpb5rBNKqWd86BX7P%20JA%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%3A4ed1756c-e31f-4809-a873-056c3026afb4%22%7D" > request
dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "g"  "WebService--AWS-BounceNotification" "Notification"


stopMojoServer
echo "RENTROLL SERVER STOPPED" 

logcheck
