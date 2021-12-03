package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/eks-anywhere-test-tool/pkg/awsprofiles"
	"github.com/aws/eks-anywhere-test-tool/pkg/constants"
	"github.com/aws/eks-anywhere-test-tool/pkg/logger"
	"os"
)

type S3 struct {
	session *session.Session
	svc     *s3.S3
}

func New(account awsprofiles.EksAccount) (*S3, error) {
	logger.V(2).Info("creating S3 client")
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: account.ProfileName(),
		Config: aws.Config{
			Region: aws.String(constants.AwsAccountRegion),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})

	if err != nil {
		fmt.Printf("Got error when setting up session: %v", err)
		os.Exit(1)
	}

	svc := s3.New(sess)
	logger.V(2).Info("created S3 client")

	return &S3{
		session: sess,
		svc:     svc,
	}, nil
}

func (s *S3) ListObjects(bucket string, prefix string) (listedObjects []*s3.Object, err error) {
	var nextToken *string
	var objects []*s3.Object

	input := &s3.ListObjectsV2Input{
		Bucket:              aws.String(bucket),
		Prefix:              aws.String(prefix),
		ContinuationToken:   nextToken,
	}

	for {
		l, err := s.svc.ListObjectsV2(input)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", err)
		}
		objects = append(objects, l.Contents...)
		if l.NextContinuationToken == nil || nextToken != nil && *nextToken == *l.NextContinuationToken {
			logger.Info("finished fetching objects", "bucket", bucket, "prefix", prefix)
			logger.V(3).Info("token comparison", "nextToken", nextToken, "nextContinuatonToken", l.NextContinuationToken)
			break
		}
		nextToken = l.NextContinuationToken
		logger.Info("fetched objects", "bucket", bucket, "prefix", prefix, "events", len(l.Contents))
		logger.V(3).Info("token comparison", "nextToken", nextToken, "nextContinuatonToken", l.NextContinuationToken)
	}
	return objects, nil
}

func (s *S3) GetObject(bucket string, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	obj, err := s.svc.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object at key %s: %v", key, err)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object at key %s: %v", key, err)
	}
	return buf.Bytes(), nil
}
