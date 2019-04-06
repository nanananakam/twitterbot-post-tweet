#!/usr/bin/env sh

aws s3 cp s3://${AWS_S3_BUCKET}/words.tar.xz .
tar Jxvf words.tar.xz
/main