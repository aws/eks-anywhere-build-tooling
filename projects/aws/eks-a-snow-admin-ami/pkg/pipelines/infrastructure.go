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
	instanceProf := p.instanceProfile()
	err = instanceProf.create(ctx, session)
	if err != nil {
		return "", err
	}

	return instanceProf.name, nil
}

func (p *Pipeline) instanceProfile() *instanceProfile {
	return &instanceProfile{
		name: fmt.Sprintf("%s-instance-profile", p.ValidNameForARN()),
		role: &role{
			name: fmt.Sprintf("%s-role", p.ValidNameForARN()),
			policyARNs: []string{
				"arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore",
				"arn:aws:iam::aws:policy/EC2InstanceProfileForImageBuilder",
			},
			assumeRolePolicyDocument: `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}`,
		},
	}
}

func isAlreadyExistIAM(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		return aerr.Code() == iam.ErrCodeEntityAlreadyExistsException
	}

	return false
}
