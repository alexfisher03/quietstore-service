package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ensureBucket(ctx context.Context, s3c *s3.Client, bucket string) {
	if _, err := s3c.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)}); err == nil {
		return
	}
	if _, err := s3c.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		if _, hb := s3c.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)}); hb == nil {
			return
		}
	}
}
