package etcdadm

func buildFlags(repository, version, cipherSuites string) []string {
	return []string{
		"-l", "debug",
		"--version", version,
		"--init-system", "kubelet",
		"--image-repository", repository,
		"--certs-dir", certDir,
		"--data-dir", dataDir,
		"--kubelet-pod-manifest-path", podSpecDir,
		"--cipher-suites", cipherSuites,
	}
}
