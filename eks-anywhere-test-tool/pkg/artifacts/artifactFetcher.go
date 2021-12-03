package artifacts

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"strings"

	"github.com/aws/eks-anywhere-test-tool/pkg/codebuild"
	"github.com/aws/eks-anywhere-test-tool/pkg/constants"
	filewriter2 "github.com/aws/eks-anywhere-test-tool/pkg/filewriter"
	"github.com/aws/eks-anywhere-test-tool/pkg/logger"
	"github.com/aws/eks-anywhere-test-tool/pkg/s3"
)

type FetchArtifactsOpt func(options *fetchArtifactConfig) (err error)

func WithCodebuildBuild(buildId string) FetchArtifactsOpt {
	return func(options *fetchArtifactConfig) (err error) {
		options.buildId = buildId
		logger.Info("user provided build ID detected", "buildId", buildId)
		return err
	}
}

type fetchArtifactConfig struct {
	buildId string
	bucket  string
}

type testArtifactFetcher struct {
	testAccountS3Client         *s3.S3
	buildAccountCodebuildClient *codebuild.Codebuild
	writer                      filewriter2.FileWriter
}

func New(testAccountS3Client *s3.S3, buildAccountCodebuildCient *codebuild.Codebuild, writer filewriter2.FileWriter) *testArtifactFetcher {
	return &testArtifactFetcher{
		testAccountS3Client:         testAccountS3Client,
		buildAccountCodebuildClient: buildAccountCodebuildCient,
		writer:                      writer,
	}
}

func (l *testArtifactFetcher) FetchArtifacts(opts ...FetchArtifactsOpt) error {
	config := &fetchArtifactConfig{
		buildId: *l.buildAccountCodebuildClient.FetchLatestBuildForProject().Id,
		bucket: os.Getenv(constants.E2eArtifactsBucketEnvVar),
	}

	for _, opt := range opts {
		err := opt(config)
		if err != nil {
			return fmt.Errorf("failed to set options on fetch artifacts config: %v", err)
		}
	}

	logger.Info("Fetching build artifacts...")

	objects, err := l.testAccountS3Client.ListObjects(config.bucket, config.buildId)
	logger.V(5).Info("Listed objects", "bucket", config.bucket, "prefix", config.buildId, "objects", len(objects))
	if err != nil {
		return fmt.Errorf("error listing objects in bucket %s at key %s: %v", config.bucket, config.buildId, err)
	}

	errs, _ := errgroup.WithContext(context.Background())

	for _, object := range objects {
		if excludedKey(*object.Key) {
			continue
		}
		obj := *object
		errs.Go(func() error {
			logger.Info("Fetching object", "key", obj.Key, "bucket", config.bucket)
			o, err := l.testAccountS3Client.GetObject(config.bucket, *obj.Key)
			if err != nil {
				return err
			}
			logger.Info("Fetched object", "key", obj.Key, "bucket", config.bucket)

			logger.Info("Writing object to file", "key", obj.Key, "bucket", config.bucket)
			err = l.writer.WriteS3KeyToFile(*obj.Key, o)
			if err != nil {
				logger.Info("error occured while writing file", "err", err)
				return fmt.Errorf("error writing object %s from bucket %s to file: %v", *obj.Key, config.bucket, err)
			}
			return nil
		})
	}
	return errs.Wait()
}

func excludedKey(key string) bool {
	excludedKeys := []string{
		"/.git/",
	}

	excludedSuffixes := []string{
		"/e2e.test",
		"/eksctl-anywhere",
	}

	for _, s := range excludedKeys {
		if strings.Contains(key, s) {
			return true
		}
	}

	for _, s := range excludedSuffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	return false
}