package pipelines

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"

	"github.com/aws/eks-anywhere-build-tooling/projects/aws/eks-a-snow-admin-ami/pkg/session"
)

type instanceProfile struct {
	name string
	role *role
}

func (i *instanceProfile) create(ctx context.Context, session *session.Session) error {
	log.Printf("Creating instance profile [%s]\n", i.name)
	iamService := iam.New(session)
	_, err := iamService.CreateInstanceProfileWithContext(ctx, &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(i.name),
	})

	if err != nil && !isAlreadyExistIAM(err) {
		return errors.Wrapf(err, "failed creating instance profile [%s]", i.name)
	}

	if err = i.addRole(ctx, session); err != nil {
		return err
	}

	return nil
}

func (i *instanceProfile) addRole(ctx context.Context, session *session.Session) error {
	log.Printf("Adding role [%s] to instance profile [%s]\n", i.role.name, i.name)
	iamService := iam.New(session)

	instanceProfile, err := iamService.GetInstanceProfileWithContext(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(i.name),
	})
	if err != nil {
		return errors.Wrapf(err, "failed checking instance profile %s for role %s", i.name, i.role.name)
	}

	// If Role already in instance profile, skip
	for _, r := range instanceProfile.InstanceProfile.Roles {
		if *r.RoleName == i.role.name {
			log.Printf("Role [%s] already in instance profile [%s]\n", i.role.name, i.name)
			return nil
		}
	}

	if err := i.role.create(ctx, session); err != nil {
		return err
	}

	_, err = iamService.AddRoleToInstanceProfileWithContext(ctx, &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(i.name),
		RoleName:            aws.String(i.role.name),
	})
	if err != nil {
		return errors.Wrapf(err, "failed adding role %s to instance profile %s", i.role.name, i.name)
	}

	return nil
}

type role struct {
	name                     string
	assumeRolePolicyDocument string
	policyARNs               []string
}

func (r *role) create(ctx context.Context, session *session.Session) error {
	log.Printf("Creating role [%s]\n", r.name)
	iamService := iam.New(session)
	_, err := iamService.CreateRoleWithContext(ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(r.name),
		Path:                     aws.String("/"),
		AssumeRolePolicyDocument: aws.String(r.assumeRolePolicyDocument),
	})

	if err != nil && !isAlreadyExistIAM(err) {
		return errors.Wrapf(err, "failed creating role %s", r.name)
	}

	for _, p := range r.policyARNs {
		log.Printf("Attaching policy [%s] role [%s]\n", p, r.name)
		_, err = iamService.AttachRolePolicyWithContext(ctx, &iam.AttachRolePolicyInput{
			RoleName:  aws.String(r.name),
			PolicyArn: aws.String(p),
		})

		if err != nil && !isAlreadyExistIAM(err) {
			return errors.Wrapf(err, "failed adding policy [%s] to role [%s]", p, r.name)
		}
	}

	return nil
}
