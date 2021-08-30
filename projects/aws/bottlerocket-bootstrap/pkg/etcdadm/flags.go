package etcdadm

func buildFlags(repository, version string) []string {
	return []string{
		"-l", "debug",
		"--version", version,
		"--init-system", "kubelet",
		"--image-repository", repository,
		"--certs-dir", certDir,
		"--data-dir", dataDir,
		"--podspec-dir", podSpecDir,
	}
}
