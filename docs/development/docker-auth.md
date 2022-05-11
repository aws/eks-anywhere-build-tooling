# Docker Registry Login

In the course of working with EKS Anywhere, we frequently need to authenticate to a Docker or OCI registry. This document aims to describe how that is accomplished, as well as optional steps that can make the process more secure and more streamlined.

Going forward, I'll use the term registry to refer to either a [Docker](https://docs.docker.com/registry/) or [OCI-compliant](https://github.com/opencontainers/distribution-spec/blob/main/spec.md) registry. In my experience, they all handle authentication in mostly the same way.

Of course Docker isn't the only tool we use when we need to authenticate to a registry, there are plenty of others: [skopeo](https://github.com/containers/skopeo), [crane](https://github.com/google/go-containerregistry), [oras-cli](https://oras.land/cli/), [helm](https://github.com/helm/helm), etc. Throughout this document, when I refer to logging into Docker, that generally means that I'm talking about any of the tools we use to interact with a registry.

## How Docker (and related tools) handle authentication

When you run `docker login $REGISTRY_HOSTNAME` for the first time, the docker command-line tool prompts you for a username and password. It uses these values to authenticate you to `$REGISTRY_HOSTNAME`. It then stores these credentials in a file on disk. The next time you interact with that registry, the docker tool references the file on disk to retrieve your username and password so you don't have to enter that information for each and every interaction with the registry. This allows scripted interactions like build processes and automated tests to run without having to pause several times for you to enter your login and password.

## Why the defaults are bad

While the default storage of credentials is convenient, it presents two problems for our use cases. The first is that the credentials are stored unencrypted. The second is that the credentials received from AWS expire after 12 hours<sup>[1](#fn-1)</sup>.

Let's start with the encryption first. On my Linux laptop, the default location for Docker's credentials is `~/.docker/config.json`. It looks like this:

```sh
$ cat ~/.docker/config.json
```
```json
{
	"auths": {
		"$REGISTRY_HOSTNAME": {
			"auth": "QVdTOmV5... <lots of redacted text>"
		}
	}
}
```

If I take that "auth" field and base64 decode it, we can see:

```sh
jq -r .auths."$REGISTRY_HOSTNAME".auth ~/.docker/config.json | base64 -d
```

```
AWS:eyJwYXls... <lots of redacted text>
```

You can see my login (AWS), and the beginning of the token (also base64 encoded) that was returned when I ran `aws ecr get-login-password`.

As for the expiration issue, the Amazon ECR documentation indicates that tokens returned by `aws ecr get-login-password` expire after 12 hours<sup>[1](#fn-1)</sup>. Unfortunately, Docker's configuration file doesn't allow any way to indicate an expiration. This can lead to cases where expired tokens are presented by Docker or related tools to Amazon ECR, which rejects them. This is often the reason why we need to `docker logout $REGISTRY_HOSTNAME`. Logging out flushes the expired token, causing Docker to prompt us for that information the next time it's needed.

## Is there a better solution?

There are some ways we can improve the situation. To tackle the lack of encryption, we can use a "Docker credential helper". A Docker credential helper is a program that when prompted for credentials for a specific registry, returns those credentials. They're programs that live in your shell's execution path, and have names that start with "docker-credential-". The specific credential helpers we're interested in are ones that interface with an encrypted secret store. There are a handful of different credential helpers that fit the bill, but it's probably easiest to use the secret store that's native to your host operating system / environment. For macOS, that's [docker-credential-osxkeychain](https://github.com/docker/docker-credential-helpers) whereas for Linux, it's usually [docker-credential-secretservice](https://github.com/docker/docker-credential-helpers).

To use a credential helper, we should first flush our existing credentials, then we can configure docker to use the helper of our choice, per the instructions linked above.

```sh
docker logout $REGISTRY_HOSTNAME
jq '.+{"credsStore":"secretservice"}' ~/.docker/config.json | sponge ~/.docker/config.json
```

And then store new ones, but this time, with encryption:

```sh
aws ecr get-login-password | docker login --username AWS --password-stdin $REGISTRY_HOSTNAME
```

Now if I look at `~/.docker/config.json`, it contains:

```json
{
  "auths": {
  },
  "credsStore": "secretservice"
}
```

Note the lack of unencrypted authentication tokens. If I list my OS's secret store, I see a new entry has been created:

```sh
secret-tool search --all server $REGISTRY_HOSTNAME
```

```
[/org/freedesktop/secrets/collection/login/22]
label = $REGISTRY_HOSTNAME
secret = eyJwYXls... <redacted>
created = 2022-05-09 17:29:28
modified = 2022-05-09 17:31:28
schema = io.docker.Credentials
attribute.username = AWS
attribute.server = $REGISTRY_HOSTNAME
attribute.docker_cli = 1
attribute.label = Docker Credentials
```

While the secret above is the same as the one found previously in `~/.docker/config.json`, the difference is that this one is stored encrypted. That's a definite plus.

Okay, so now we have our Docker credentials encrypted, but what about their expiration? To solve this problem, we look to the [Amazon ECR Docker Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper). This is yet another docker credential helper, but this one's configured to interact with your AWS command-line tool's credentials<sup>[2](#fn-2)</sup> to refresh your tokens as needed. Once you have `docker-credential-ecr-login` installed and configured, you won't have to enter your password at all! Whenever docker detects the credentials are out of date, the Amazon ECR Docker Credential Helper fetches a new token and supplies it to Docker for you.

At this point, if you're following along, your Docker configuration file should look something like this:

```json
{
	"auths": {
	},
	"credsStore": "secretservice",
	"credHelpers": {
		"$REGISTRY_HOSTNAME": "ecr-login",
	}
}
```

## Encrypt all the things

With this configuration, Docker credentials are encrypted, and Amazon ECR tokens are refreshed automatically. Life is good. But there's one more thing we can do to improve our security. Some of you might be asking yourselves, "Hey, this Amazon ECR docker credential helper is great, but just how does it perform its authentication?" I'm glad you asked. If you followed the [installation and configuration](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) for the `aws` CLI tool, you ran `aws configure`, which prompted you for your AWS Access Key ID and your AWS Secret Access Key. It then stored those... You guessed it! Unencrypted in a file on your disk.

Fortunately, there's a solution for this as well, and it involves using our good friend the secret store again. The [AWS Vault](https://github.com/99designs/aws-vault) project provides a way to store and access your encrypted AWS credentials via OS-native secret store plus a few other methods as well. Once that is installed and configured, you can delete your `~/.aws/credentials` file---your credentials are stored encrypted, and available via the `aws-vault` command-line:

```sh
# One-time configuration for your $SHELL of choice (bash shown here):
export AWS_VAULT_BACKEND=secret-service
# See all of the env vars at:
# https://github.com/99designs/aws-vault/blob/master/USAGE.md#environment-variables
aws-vault exec default -- skopeo inspect docker://$REGISTRY_HOSTNAME/eks-anywhere-packages@sha256:3a189f892778cbef296f38656bc1b1b0fa4100107c6dfb75598ca1d45b1d4172
{
    "Name": "$REGISTRY_HOSTNAME/eks-anywhere-packages",
    "Digest": "sha256:3a189f892778cbef296f38656bc1b1b0fa4100107c6dfb75598ca1d45b1d4172",
    "RepoTags": [
        "v0.1.6-c7cade6f65080958cb790743bd169e930ac27b58",
        "latest",
        "0.1.6-c7cade6f65080958cb790743bd169e930ac27b58"
    ],
    "Created": "2022-05-05T10:58:37.545569852-06:00",
    "DockerVersion": "",
    "Labels": null,
    "Architecture": "amd64",
    "Os": "linux",
    "Layers": [
        "sha256:6dc22987e65918b540d1e5f980d351ede1acf76f3168d3166457b289a2f7101b",
        "sha256:7445ad857b5293f6b92bfc2720db3bee0cc94aadf9c9c8f1e62d1211e7cd5b18",
        "sha256:b83f65c9a9f1f9a0d9f55546ab33efa6ffebabbfe628ad777f61df39a2b72fc6",
        "sha256:2534fb3dcf18bccc2cb94e1de8af861ebfc241d517c24693f2e0bbe17811369d"
    ],
    "Env": [
        "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    ]
}
```

## The end?

With that, or journey is complete. Or nearly so. It's possible that you're wondering, "Why did we install the Docker credential login helper at all? We installed the Amazon ECR Docker Credential Helper right after, and they perform the same jobs." You caught me; we didn't need the Docker credential login helper. However, before you get mad at me, please understand that I had your best interests at 🤍. The reason these instructions started with the Docker credential login helper was so that any non-AWS registries you might authenticate to will gain the benefit of having its credentials stored via your encrypted secret store. Improved security for everyone!

## Appendix A: Where Docker (and friends) store credentials

In my adventures, I've found a great number of places that Docker and friends can read or write credentials. I doubt this is all of them, but these are the ones I've found so far. I'm listing them here to expose their existence in the hopes that it won't take you hours of debugging to find them if you ever have authentication problems with some random Docker-related tool.

For each file below, I've included the source where I found the information. The notation `some-term(5)` means "section 5 of the some-term man page" and can be accessed from your shell via `man 5 some-term`.

```
    # Docker's default on Linux:
	# https://docs.docker.com/engine/reference/commandline/cli/#environment-variables
    ${DOCKER_CONFIG:-${HOME}/.docker}/config.json

    # Helm https://helm.sh/docs/helm/helm_env/
    ${HELM_REGISTRY_CONFIG:-${XDG_CONFIG_HOME:-${HOME}/.config}}/helm/registry/config.json
	# Also be aware of the value of $HELM_CONFIG_HOME

    # containers-auth.json(5) (Linux hosts)
    ${XDG_RUNTIME_DIR}/containers/auth.json

    # containers-auth.json(5) (macOS and Windows hosts)
    ${XDG_CONFIG_HOME:-${HOME}/.config}/containers/auth.json

    # Why is Linux different above? I'm not sure, but I'd guess that
    # macOS and Windows don't follow the XDG Base Directory spec, and
    # so don't have anything equivalent to $XDG_RUNTIME_DIR, and so
    # they chose the next best place, which was $XDG_CONFIG_HOME,
    # which probably _also_ doesn't exist outside of Linux, but
    # specifies the use of $HOME as a fallback, and that's something
    # that both Windows and macOS do have.

    # skopeo-login(1)
    ${REGISTRY_AUTH_FILE}

    # A fallback
    # https://github.com/docker/cli/blob/2291f610ae73533e6e0749d4ef1e360149b1e46b/cli/config/config.go#L111-L123
    # and also containers-auth.json(5)
    ${HOME}/.dockercfg
```

I recommend setting all of the applicable environment variables to all refer to one file. I chose my Docker config file. In theory this means that skopeo, crane, helm, oras-cli, et al will also use my configured credentials helpers.

Some tools will check for other tools' files. For example, the oras CLI tool will read your Docker configuration file if it exists. The files that the various tools look for, and in what order can vary widely.

## Appendix B: An "av" shell function

Here's a function I've found handy. I've added this to my ~/.bash_profile:

```sh
# shellcheck disable=SC2219
av () {
    local buf
    local nArgs=0
    local profile

    for arg in "$@"; do
		if [ "$arg" = "--" ]; then
		    let nArgs=nArgs+1
		    shift $nArgs
		    profile=$buf
		    break
		fi
		buf="${buf}${buf:+ }$arg"
		let nArgs=nArgs+1
    done
    aws-vault exec "${profile:-default}" -- "$@"
}
export -f av
```

Then I can run:

    ```sh
	av make images
	```

To use my default profile, or:

	```sh
	av other-profile make images
	```

To use a different profile.

## Footnotes

<a name="fn-1"></a>
1. https://docs.aws.amazon.com/AmazonECR/latest/userguide/registry_auth.html#registry-auth-token

<a name="fn-2"></a>
2. Assuming you set them up via `aws configure` as part of your AWS CLI installation.
