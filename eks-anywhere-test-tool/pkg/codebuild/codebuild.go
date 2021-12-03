package codebuild

import (
	"fmt"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"

	"github.com/aws/eks-anywhere-test-tool/pkg/awsprofiles"
	"github.com/aws/eks-anywhere-test-tool/pkg/constants"
	"github.com/aws/eks-anywhere-test-tool/pkg/logger"
)

type Codebuild struct {
	session *session.Session
	svc     *codebuild.CodeBuild
}

func New(account awsprofiles.EksAccount) (*Codebuild, error) {
	logger.V(2).Info("creating codebuild client")
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

	svc := codebuild.New(sess)
	logger.V(2).Info("created codebuild client")

	return &Codebuild{
		session: sess,
		svc:     svc,
	}, nil
}

func (c *Codebuild) FetchBuildForProject(id string) *codebuild.Build {
	return c.getBuildById(id)
}

func (c *Codebuild) FetchLatestBuildForProject() *codebuild.Build {
	builds := c.FetchBuildsForProject()
	latestId := *builds.Ids[0]
	return c.getBuildById(latestId)
}

func (c *Codebuild) getBuildById(id string) *codebuild.Build {
	i := []*string{aws.String(id)}
	latestBuild, err := c.svc.BatchGetBuilds(&codebuild.BatchGetBuildsInput{Ids: i})
	if err != nil {
		fmt.Printf("Got an error when fetching latest build for project: %v", err)
		os.Exit(1)
	}
	return latestBuild.Builds[0]
}

func (c *Codebuild) FetchBuildsForProject() *codebuild.ListBuildsForProjectOutput {
	builds, err := c.svc.ListBuildsForProject(&codebuild.ListBuildsForProjectInput{
		NextToken:   nil,
		ProjectName: aws.String(constants.EksATestCodebuildProject),
		SortOrder:   aws.String(codebuild.SortOrderTypeDescending),
	})
	if err != nil {
		fmt.Printf("Got an error when fetching builds for project: %v", err)
		os.Exit(1)
	}
	return builds
}
