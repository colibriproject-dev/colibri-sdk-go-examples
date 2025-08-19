#!/bin/bash

awslocal sns create-topic --name SCHOOL_ENROLLMENT_CREATED
awslocal sqs create-queue --queue-name SCHOOL_ENROLLMENT_CREATED_TOPIC_TEST
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:SCHOOL_ENROLLMENT_CREATED \
    --protocol sqs \
    --notification-endpoint arn:aws:sqs:us-east-1:queue:SCHOOL_ENROLLMENT_CREATED_TOPIC_TEST

awslocal sns create-topic --name SCHOOL_ENROLLMENT_DELETED
awslocal sqs create-queue --queue-name SCHOOL_ENROLLMENT_DELETED_TOPIC_TEST
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:SCHOOL_ENROLLMENT_DELETED \
    --protocol sqs \
    --notification-endpoint arn:aws:sqs:us-east-1:queue:SCHOOL_ENROLLMENT_DELETED_TOPIC_TEST

awslocal sns create-topic --name SCHOOL_COURSE_DELETED
awslocal sqs create-queue --queue-name SCHOOL_COURSE_DELETED_TOPIC_TEST
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:SCHOOL_COURSE_DELETED \
    --protocol sqs \
    --notification-endpoint arn:aws:sqs:us-east-1:queue:SCHOOL_COURSE_DELETED_TOPIC_TEST

awslocal sns create-topic --name SCHOOL_STUDENT_DELETED
awslocal sqs create-queue --queue-name SCHOOL_STUDENT_DELETED_TOPIC_TEST
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:SCHOOL_STUDENT_DELETED \
    --protocol sqs \
    --notification-endpoint arn:aws:sqs:us-east-1:queue:SCHOOL_STUDENT_DELETED_TOPIC_TEST

awslocal sns create-topic --name FINANCIAL_INSTALLMENT
awslocal sqs create-queue --queue-name FINANCIAL_INSTALLMENT_SCHOOL
awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:FINANCIAL_INSTALLMENT \
    --protocol sqs \
    --notification-endpoint arn:aws:sqs:us-east-1:queue:FINANCIAL_INSTALLMENT_SCHOOL

awslocal s3api create-bucket --bucket meu-bucket --acl public-read

echo "localstack topics and queues started"