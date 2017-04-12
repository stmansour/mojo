package main

// This module handles the bounce messages from AWS SNS.
//
// For reference:
// 	bounce object: http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html#bounce-object
// 	complaint object: http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html#complaint-object
// 	delivery object: http://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html#delivery-object
