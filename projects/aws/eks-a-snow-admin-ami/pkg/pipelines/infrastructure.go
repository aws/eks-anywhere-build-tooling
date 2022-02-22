package pipelines

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

type infraConfig struct {
	Name string
}

func (p *Pipeline) setupInfraConfig(ctx context.Context, session *session.Session) (arn string, err error) {
	log.Printf("Setting up infra config for pipeline [%s]\n", p.Name)
	instanceProfileName, err := p.setupInstanceProfile(ctx, session)
	if err != nil {
		return "", err
	}

	builder := imagebuilder.New(session)
	name := fmt.Sprintf("%s-infra-config", p.ValidNameForARN())

	log.Printf("Creating infra config with instance profile name [%s] for pipeline [%s]\n", instanceProfileName, p.Name)
	_, err = builder.CreateInfrastructureConfigurationWithContext(ctx, &imagebuilder.CreateInfrastructureConfigurationInput{
		Name:                aws.String(name),
		InstanceProfileName: &instanceProfileName,
		InstanceTypes:       aws.StringSlice([]string{"m4.2xlarge"}),
	})

	if err != nil && !isAlreadyExist(err) {
		return "", errors.Wrapf(err, "failed creating infra config for pipeline %s", p.Name)
	}

	log.Printf("Searching for infra config [%s] for pipeline [%s]\n", name, p.Name)
	infraConfigList, err := builder.ListInfrastructureConfigurationsWithContext(ctx, &imagebuilder.ListInfrastructureConfigurationsInput{
		Filters: []*imagebuilder.Filter{
			{
				Name:   aws.String("name"),
				Values: aws.StringSlice([]string{name}),
			},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed searching for infra config [%s] ARN for pipeline [%s]", name, p.Name)
	}

	if len(infraConfigList.InfrastructureConfigurationSummaryList) == 0 {
		return "", errors.Errorf("could not find infra config [%s] %v", name, infraConfigList)
	}

	return *infraConfigList.InfrastructureConfigurationSummaryList[0].Arn, nil
}

func (p *Pipeline) setupInstanceProfile(ctx context.Context, session *session.Session) (name string, err error) {
	instanceProfileName, err := p.createInstanceProfile(ctx, session)
	if err != nil {
		return "", err
	}

	return instanceProfileName, nil
}

func (p *Pipeline) createInstanceProfile(ctx context.Context, session *session.Session) (name string, err error) {
	name = fmt.Sprintf("%s-instance-profile", p.ValidNameForARN())
	log.Printf("Creating instance profile [%s] for pipeline [%s]\n", name, p.Name)
	iamService := iam.New(session)
	_, err = iamService.CreateInstanceProfileWithContext(ctx, &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	})

	if isAlreadyExistIAM(err) {
		return name, nil
	}

	if err != nil {
		return "", errors.Wrapf(err, "failed creating instance profile for pipeline %s", p.Name)
	}

	if err = addRoleToInstanceProfile(ctx, session, name, "EC2InstanceProfileForImageBuilder"); err != nil {
		return "", err
	}

	return name, nil
}

func addRoleToInstanceProfile(ctx context.Context, session *session.Session, instanceProfileName, roleName string) error {
	log.Printf("Adding role [%s] to instance profile [%s]\n", roleName, instanceProfileName)
	iamService := iam.New(session)

	_, err := iamService.AddRoleToInstanceProfileWithContext(ctx, &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(instanceProfileName),
		RoleName:            aws.String(roleName),
	})
	if err != nil {
		return errors.Wrapf(err, "failed adding role %s to instance profile %s", roleName, instanceProfileName)
	}

	return nil
}

func isAlreadyExistIAM(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code() == iam.ErrCodeEntityAlreadyExistsException
	}

	return false
}
