#!/bin/bash

awslocal s3 mb s3://ecommerce-uploads
awslocal sqs create-queue --queue-name ecommerce-events

echo "LocalStach initialization complete"
