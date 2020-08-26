#!/bin/bash

#---------------------------------------------------------------
# TOP is the directory where RentRoll begins. It is used
# in base.sh to set other useful directories such as ${BASHDIR}
#---------------------------------------------------------------
TOP=../..
BINDIR=${TOP}/tmp/mojo

TESTNAME="Web Services"
TESTSUMMARY="Test Web Services"
SMALLDB="smalldb.sql"
FAASCRAPE="${BINDIR}/scrapefaa"
NEWDB="${BINDIR}/mojonewdb"

CREATENEWDB=0

#---------------------------------------------------------------
#  Use the testdb for these tests...
#---------------------------------------------------------------
echo "Create new database..."

if [ ! -f "${SMALLDB}" ]; then
	${NEWDB}
	${FAASCRAPE} -q
	mysqldump --no-defaults mojo >${SMALLDB}
fi

mysql --no-defaults mojo < ${SMALLDB}

source ../share/base.sh

echo "STARTING MOJO SERVER"
startMojoServer

curl http://localhost:8275/v1/ping/


#------------------------------------------------------------------------------
#  TEST a
#  basic tests
#
#  Scenario:
#  follow the comments below
#
#  Expected Results:
#   ...
#------------------------------------------------------------------------------
TFILES="a"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
    echo "Test ${TFILES}"
	echo "request%3D%7B%22cmd%22%3A%22get%22%2C%22selected%22%3A%5B%5D%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
	dojsonPOST "http://localhost:8275/v1/people/1" "request" "${TFILES}${STEP}"  "WebService--GetPeople"
	dojsonGET "http://localhost:8275/v1/peoplecount/" "${TFILES}${STEP}"  "WebService--PeopleCount"
	dojsonGET "http://localhost:8275/v1/groupcount/" "${TFILES}${STEP}"  "WebService--GroupCount"
	dojsonGET "http://localhost:8275/v1/peoplestats/" "${TFILES}${STEP}"  "WebService--PeopleStats"

	echo "request%3D%7B%22cmd%22%3A%22%22%2C%22limit%22%3A100%2C%22offset%22%3A0%7D" > request
	dojsonPOST "http://localhost:8275/v1/queries/" "request" "${TFILES}${STEP}"  "WebService--SearchQueries"

	dojsonGET "http://localhost:8275/v1/qrescount/1" "${TFILES}${STEP}"  "WebService--QueryResultsCount"

	# Subscription Succeeded Notification
	echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%22b7001127-9d1f-524a-9045-404e3913897b%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%22%2C%22Message%22%3A%7B%22notificationType%22%3A%22AmazonSnsSubscriptionSucceeded%22%2C%22message%22%3A%22You%20have%20successfully%20subscribed%20your%20Amazon%20SNS%20topic%20'arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email'%20to%20receive%20'Bounce'%20notifications%20from%20Amazon%20SES%20for%20identity%20'sman%40accordinterests.com'.%22%7D%2C%22Timestamp%22%3A%222017-04-13T16%3A05%3A32.839Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22Dv6o%20yiGfg3dzDkRTZ2Wf4Uj08DtzhziuPW99QkGDes9pLd3lOIUwl%20u5XEfURI374NiZD3pr7Ku7U7WwVCA%20Laa4RlGF8T5FXlpaLuDYgccEgoDNm1IAX%204a19yUdWD7dAzv1eEIcQSTgezT4QGYZh1XuTAvUAU%2FZnBcfL5P8nGAAKrf78jEUkV0rcH%2FHRw1ndewlzqcQoDy2H64x%20EqHuYIO%2F7mbTgDFQNxaQg0Fr64ipzheVfijvWe4NpWjblAAawuJRnPyEpDfFdJZFbfk%206y8QpC8O56QDVXV5h5PaKvioqo0KyOyK9HpuTGRgy%20X5OZAxjKtLXxwVZF%20wzdQ%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%3A4ed1756c-e31f-4809-a873-056c3026afb4%22%7D" > request
	dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "${TFILES}${STEP}"  "WebService--AWS-SubscriptionSucceeded" "Notification"

	# Bounce Notification
	echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%2239a64847-8102-5692-bcd7-192d924f075d%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%22%2C%22Message%22%3A%22%7B%5C%22notificationType%5C%22%3A%5C%22Bounce%5C%22%2C%5C%22bounce%5C%22%3A%7B%5C%22bounceType%5C%22%3A%5C%22Permanent%5C%22%2C%5C%22bounceSubType%5C%22%3A%5C%22General%5C%22%2C%5C%22bouncedRecipients%5C%22%3A%5B%7B%5C%22emailAddress%5C%22%3A%5C%22bounce%40simulator.amazonses.com%5C%22%2C%5C%22action%5C%22%3A%5C%22failed%5C%22%2C%5C%22status%5C%22%3A%5C%225.1.1%5C%22%2C%5C%22diagnosticCode%5C%22%3A%5C%22smtp%3B%20550%205.1.1%20user%20unknown%5C%22%7D%5D%2C%5C%22timestamp%5C%22%3A%5C%222017-04-14T00%3A00%3A13.068Z%5C%22%2C%5C%22feedbackId%5C%22%3A%5C%220100015b69c29ab8-fbf8d711-ea80-4465-bb85-aa901ea7f145-000000%5C%22%2C%5C%22remoteMtaIp%5C%22%3A%5C%22207.171.163.188%5C%22%2C%5C%22reportingMTA%5C%22%3A%5C%22dsn%3B%20a8-50.smtp-out.amazonses.com%5C%22%7D%2C%5C%22mail%5C%22%3A%7B%5C%22timestamp%5C%22%3A%5C%222017-04-14T00%3A00%3A08.000Z%5C%22%2C%5C%22source%5C%22%3A%5C%22sman%40stevemansour.com%5C%22%2C%5C%22sourceArn%5C%22%3A%5C%22arn%3Aaws%3Ases%3Aus-east-1%3A553053000801%3Aidentity%2Fsman%40stevemansour.com%5C%22%2C%5C%22sourceIp%5C%22%3A%5C%2224.6.191.18%5C%22%2C%5C%22sendingAccountId%5C%22%3A%5C%22553053000801%5C%22%2C%5C%22messageId%5C%22%3A%5C%220100015b69c28a6f-956eba47-89be-4ea8-8e5e-9924ee3954eb-000000%5C%22%2C%5C%22destination%5C%22%3A%5B%5C%22bounce%40simulator.amazonses.com%5C%22%5D%7D%7D%22%2C%22Timestamp%22%3A%222017-04-14T00%3A00%3A13.118Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22Akcg6IxxVQ1aIC0YIJMnuY9FVaNpndIQjuzVVFWVHZ7n1jKjIZYeH%2Bt1AAQBh%2FXcucIVlhKKfKjqgLTimPpde%2B4FEc6JS6G3eQW1pxkKuyDlPdj0AjUbMouCKNKldmFGWlA%2BuAF2GzXHUQPAwPMpnl%2FcyIySJgqlQk52ODfkzdJ1EFQM6D2jKPXrs1B5OdSgb9e3uGCo4R3U3wQkFbEVtj7%2BTbsv9AAEMwrVyrhhUDlBitFMpBVkI0Cu%2BVGf6V%2BefQtBe40l8t%2BEHVzdfDy2RYDH57PgnOVqENRR3NGsZ2ugqdKlLfJaVEgfhug6fi%2F4i%2FmQndXNhGLlgLzz8Fzl7Q%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%3A41abeec0-7f87-419f-a33e-fd97eda671bd%22%7D" > request
	dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "${TFILES}${STEP}"  "WebService--AWS-BounceNotification" "Notification"

	# Complaint Notification
	echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%222b870f02-7a74-5a86-b804-5e6862fd2b6d%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_complaint_email%22%2C%22Message%22%3A%22%7B%5C%22notificationType%5C%22%3A%5C%22Complaint%5C%22%2C%5C%22complaint%5C%22%3A%7B%5C%22complainedRecipients%5C%22%3A%5B%7B%5C%22emailAddress%5C%22%3A%5C%22complaint%40simulator.amazonses.com%5C%22%7D%5D%2C%5C%22timestamp%5C%22%3A%5C%222017-04-14T00%3A22%3A33.000Z%5C%22%2C%5C%22feedbackId%5C%22%3A%5C%220100015b69d70fea-7457e069-20a8-11e7-a6dc-412614a78c65-000000%5C%22%2C%5C%22userAgent%5C%22%3A%5C%22Amazon%20SES%20Mailbox%20Simulator%5C%22%2C%5C%22complaintFeedbackType%5C%22%3A%5C%22abuse%5C%22%7D%2C%5C%22mail%5C%22%3A%7B%5C%22timestamp%5C%22%3A%5C%222017-04-14T00%3A23%3A04.000Z%5C%22%2C%5C%22source%5C%22%3A%5C%22sman%40stevemansour.com%5C%22%2C%5C%22sourceArn%5C%22%3A%5C%22arn%3Aaws%3Ases%3Aus-east-1%3A553053000801%3Aidentity%2Fsman%40stevemansour.com%5C%22%2C%5C%22sourceIp%5C%22%3A%5C%2224.6.191.18%5C%22%2C%5C%22sendingAccountId%5C%22%3A%5C%22553053000801%5C%22%2C%5C%22messageId%5C%22%3A%5C%220100015b69d70daf-f7dec14b-7e04-4aea-9f6c-2f2026750393-000000%5C%22%2C%5C%22destination%5C%22%3A%5B%5C%22complaint%40simulator.amazonses.com%5C%22%5D%7D%7D%22%2C%22Timestamp%22%3A%222017-04-14T00%3A22%3A33.819Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22Goc0TAcgZa%2BsAVjHIx1NjeWHf2mtse2lrqgNkwktqDhQkafYvGChe1XeEl5rxsQqHMwy2YU3oNFC1Cc8QubKqemw1LavtR3K%2BZaOEYM8z%2F3Y6mn6RQQaqSb5uHaMyV2f%2FY5LUiA9Uc%2BMch565PeC4Xp6XZSjwoh2y1rxGk5cSxbr3Ln9X4I3TX4opfxK5TgLSlwnnTH15NUyhjgXV3MNvfV2MeE2bBHTqSV%2FOoK1wi1fUqIsz0vKwJCTzaVoX9IjLFcntyHAgD%2F%2FjgM5UoVSCWS1dJtkoMJT%2BPvyTmuZsrLWpfMwV7%2B%2F8l3nmHhslWJdlEeCLIWrEk%2FoI7bDBmDFdg%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_complaint_email%3A095941ef-a25c-42bb-880a-13b9f72eda2c%22%7D" > request
	dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "${TFILES}${STEP}"  "WebService--AWS-BounceNotification" "Notification"

	# Suppression List Notification
	echo "%7B%22Type%22%3A%22Notification%22%2C%22MessageId%22%3A%22b37cb42c-b023-53e1-84ab-e54794fb73d9%22%2C%22TopicArn%22%3A%22arn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%22%2C%22Message%22%3A%22%7B%5C%22notificationType%5C%22%3A%5C%22Bounce%5C%22%2C%5C%22bounce%5C%22%3A%7B%5C%22bounceType%5C%22%3A%5C%22Permanent%5C%22%2C%5C%22bounceSubType%5C%22%3A%5C%22Suppressed%5C%22%2C%5C%22bouncedRecipients%5C%22%3A%5B%7B%5C%22emailAddress%5C%22%3A%5C%22suppressionlist%40simulator.amazonses.com%5C%22%2C%5C%22action%5C%22%3A%5C%22failed%5C%22%2C%5C%22status%5C%22%3A%5C%225.1.1%5C%22%2C%5C%22diagnosticCode%5C%22%3A%5C%22Amazon%20SES%20has%20suppressed%20sending%20to%20this%20address%20because%20it%20has%20a%20recent%20history%20of%20bouncing%20as%20an%20invalid%20address.%20For%20more%20information%20about%20how%20to%20remove%20an%20address%20from%20the%20suppression%20list%2C%20see%20the%20Amazon%20SES%20Developer%20Guide%3A%20http%3A%2F%2Fdocs.aws.amazon.com%2Fses%2Flatest%2FDeveloperGuide%2Fremove-from-suppressionlist.html%20%5C%22%7D%5D%2C%5C%22timestamp%5C%22%3A%5C%222017-04-14T19%3A27%3A20.192Z%5C%22%2C%5C%22feedbackId%5C%22%3A%5C%220100015b6def2245-60a4e2d6-2148-11e7-af16-351a27feef70-000000%5C%22%2C%5C%22reportingMTA%5C%22%3A%5C%22dns%3B%20amazonses.com%5C%22%7D%2C%5C%22mail%5C%22%3A%7B%5C%22timestamp%5C%22%3A%5C%222017-04-14T19%3A27%3A19.000Z%5C%22%2C%5C%22source%5C%22%3A%5C%22sman%40stevemansour.com%5C%22%2C%5C%22sourceArn%5C%22%3A%5C%22arn%3Aaws%3Ases%3Aus-east-1%3A553053000801%3Aidentity%2Fsman%40stevemansour.com%5C%22%2C%5C%22sourceIp%5C%22%3A%5C%2224.6.191.18%5C%22%2C%5C%22sendingAccountId%5C%22%3A%5C%22553053000801%5C%22%2C%5C%22messageId%5C%22%3A%5C%220100015b6def2076-f0e2ac3c-2272-409f-85ba-bc55e4b0fe2b-000000%5C%22%2C%5C%22destination%5C%22%3A%5B%5C%22suppressionlist%40simulator.amazonses.com%5C%22%5D%7D%7D%22%2C%22Timestamp%22%3A%222017-04-14T19%3A27%3A20.255Z%22%2C%22SignatureVersion%22%3A%221%22%2C%22Signature%22%3A%22DZkZFIPvG7U8btJXlHNZ3ry89bQE7HET98BUTcUaX3yS7TR1DKmybD6g2vJ0sJypwN8xb9zJfGe23xYGQ0wB5mvLRcfw70F8gHX0VcQMswRzQXVGxr7GHKF5aEkDjYef6370ofSfPhb9j0%2FQSA2IKP8FWQe66PuJTykSuwRLB3aghbWkCgeCdr49xEaN8QnUnpS1q4URkNJSO8Tdz6TaWQjLd%209IOsqCuD%2F4HnMCSTPd2%201yJdq9lqt5hmIQOIMPRNKGgCFt9%20EDwJ%20J7p3e8mp7Q8U%2FA8NVp78xPptz832u7cp7gHrianqxtlLU6szS2fymAoKmGD7tmXKQ%203nthA%3D%3D%22%2C%22SigningCertURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2FSimpleNotificationService-b95095beb82e8f6a046b3aafc7f4149a.pem%22%2C%22UnsubscribeURL%22%3A%22https%3A%2F%2Fsns.us-east-1.amazonaws.com%2F%3FAction%3DUnsubscribe%26SubscriptionArn%3Darn%3Aaws%3Asns%3Aus-east-1%3A553053000801%3Aaccord_mojo_bounced_email%3A41abeec0-7f87-419f-a33e-fd97eda671bd%22%7D" > request
	dojsonAwsPOST "http://localhost:8275/v1/aws" "request" "${TFILES}${STEP}"  "WebService--AWS-SuppressionList" "Notification"

	# Optout FAIL
	doHtmlGET "http://localhost:8275/v1/optout?e=sman@accordinterests.com&c=ad7a567ae84ea687b23f74c5e18f0ce" "${TFILES}${STEP}"  "WebService--OptOut-fail"

	# Optout SUCCESS
	doHtmlGET "http://localhost:8275/v1/optout?e=sman@accordinterests.com&c=9010fb3ff00db43ebd891bc39d8dafc8" "${TFILES}${STEP}"  "WebService--OptOut-success"
fi

#------------------------------------------------------------------------------
#  TEST b
#
#  Read the groups associated witht the supplied PID
#
#  /v1/pgroup/UID
#
#  Scenario:
#  Search
#
#  Expected Results:
#   1.
#   2.
#------------------------------------------------------------------------------
TFILES="b"
STEP=0
if [ "${SINGLETEST}${TFILES}" = "${TFILES}" -o "${SINGLETEST}${TFILES}" = "${TFILES}${TFILES}" ]; then
	# b0
    encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8275/v1/pgroup/12" "request" "${TFILES}${STEP}"  "readGroupMembership"

	# b1
	encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8275/v1/groups/" "request" "${TFILES}${STEP}"  "groups"

	# b2: Add Steve to FAA
	encodeRequest '{"cmd":"save","Groups":[4,3,1,2,5]}'
    dojsonPOST "http://localhost:8275/v1/groupmembership/12" "request" "${TFILES}${STEP}"  "setGroupMembership-add-GID-1"

	# b3: Read Steve's group memberships and make sure that it contains GID 1
	encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8275/v1/pgroup/12" "request" "${TFILES}${STEP}"  "readGroupMembership"

	# b4: Remove Steve from FAA
	encodeRequest '{"cmd":"save","Groups":[4,3,2,5]}'
    dojsonPOST "http://localhost:8275/v1/groupmembership/12" "request" "${TFILES}${STEP}"  "setGroupMembership-remove-GID-1"

	# b5: Read Steve's group memberships and make sure that it does not contain GID 1
	encodeRequest '{"cmd":"get","selected":[],"limit":100,"offset":0}'
    dojsonPOST "http://localhost:8275/v1/pgroup/12" "request" "${TFILES}${STEP}"  "personGroupList"

	# b6: Test Transactant Typedown
    dojsonGET "http://localhost:8275/v1/grouptd/?request%3D%7B%22search%22%3A%22te%22%2C%22max%22%3A250%7D" "${TFILES}${STEP}" "GroupTypedown"

	# b7: NEW EMAIL bad group:  Test a post to add new email/name as member of group that doesn't exist
	encodeRequest '{"cmd":"save","name":"Sally Smith","email": "sally@smith.com","group":"smanmusic" }'
    dojsonPOST "http://localhost:8275/v1/addtogroup" "request" "${TFILES}${STEP}"  "addtogroup-New_Email_Bad_Group"

	# b8: NEW EMAIL:  Test a post to add new email/name as member of an existing group
	encodeRequest '{"cmd":"save","name":"Sally Smith","email": "sally@smith.com","group":"FAA" }'
    dojsonPOST "http://localhost:8275/v1/addtogroup" "request" "${TFILES}${STEP}"  "addtogroup_New_Person_NORMAL"

	# b9: BAD EMAIL ADDR:  Test a post to add new email/name as member of group
	encodeRequest '{"cmd":"save","name":"Sally Smith","email": "sally@smith","group":"smanmusic" }'
    dojsonPOST "http://localhost:8275/v1/addtogroup" "request" "${TFILES}${STEP}"  "addtogroup_InvalidEmailAddress"

	# b10: KNOWN EMAIL, NEW GROUP:
	encodeRequest '{"cmd":"save","name":"Sally Smith","email": "sally@smith.com","group":"MojoTest" }'
    dojsonPOST "http://localhost:8275/v1/addtogroup" "request" "${TFILES}${STEP}"  "addtogroup_KnownPerson_NewGroup"

	# b11: KNOWN EMAIL, ALREADY GROUP MEMBER:
	encodeRequest '{"cmd":"save","name":"Sally Smith","email": "sally@smith.com","group":"MojoTest" }'
    dojsonPOST "http://localhost:8275/v1/addtogroup" "request" "${TFILES}${STEP}"  "addtogroup_AlreadyAMember"

	# b12: MISSING EMAIL
	encodeRequest '{"cmd":"save","name":"Sally Smith","email": "","group":"MojoTest" }'
    dojsonPOST "http://localhost:8275/v1/addtogroup" "request" "${TFILES}${STEP}"  "addtogroup_MissingEmail"

	# b13: MISSING GROUP
	encodeRequest '{"cmd":"save","name":"Sally Smith","email": "sally@smith.com","group":"" }'
    dojsonPOST "http://localhost:8275/v1/addtogroup" "request" "${TFILES}${STEP}"  "addtogroup_MissingGroup"

fi


stopMojoServer
echo "RENTROLL SERVER STOPPED"

logcheck
