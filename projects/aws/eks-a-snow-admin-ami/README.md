# Snow EKS-A Admin AMI
This project builds an AMI that embeds EKS-A and all its dependencies so it can be used to spin up EKS-A clusters with the Snow provider in disconnected environments.

## How to use
### Setup
The first command sets up the basic AWS infrastructure to build the ami:
```shell
./eks-a-snow-admin-ami setup
```

The binary used the Go SDK for AWS which by default will use the same configuration as you local `aws` cli. Make sure it's configured with the right access and pointing to the right account.

### Build
You can build AMIs can be built with a different version of EKS-A and pointing to different release manifest URLs (think dev/staging/production). For example:
```shell
./eks-a-snow-admin-ami build --eksa-version v0.7.0 --eksa-release-manifest-url https://anywhere-assets.eks.amazonaws.com/releases/eks-a/manifest.yaml
```

## Testing
You can run the setup and build commands for the most complete E2E testing. However, building an AMI through Image Builder can take sometimes a long time.
When making changes to the components, you might want a quicker way to iterate.

In that case you can use the test script. This will spinup an EC2 instance, copy the necessary files and run the build and validate phases using `awstoe`.
The only two arguments needed are the name of one of your keys (however they are named in your AWS account) and the path to such key in your local disk.
```shell
./test/testComponents.sh my-key-pair ~/.ssh/my-key-pair.pem
```

You can optionally pass arguments to the components:
```shell
./test/testComponents.sh my-key-pair ~/.ssh/my-key-pair.pem EksAnywhereVersion=v0.7.2
```

The script will copy all the logs to you local disk after the test is finished in the `logs` folder. Beware that it won't terminate the ec2 instance if the test fails so you will need to terminate it manually.

The script already sets sane defaults but here are some things you can configure with env vars:
```shell
INSTANCE_TYPE # EC2 instance type to run awstoe on. Default is 't2.large'
AMI_ID # AMI to use in the test instance. Default is 'ami-0892d3c7ee96c0bf7' (Ubuntu 20 AMD64)
USER # user to SSH into the instance. It will depend on your OS. Default is 'ubuntu'
DOCUMENTS # Components to pass to awstoe (comma separated). This should default to all the components committed to the repo
PHASES # Phases to run in awstoe (comma separated). Defaults to build and validate
```

## Adding new components
Add your new yaml file to `components`. In addition, you will need to edit:
* `pkg/snow/config.go`: this file acts as the "declaration" on the aws infra. Add a reference to your new yaml file with proper name and description. If it need parameters, you will need to pass them from the `build` command.
* `test/testComponents.sh`: add your new yaml document to `DOCUMENTS`

If you are testing the binary in your own account, you will need to run the `setup` command again to setup your new component.