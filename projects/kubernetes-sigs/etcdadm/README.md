## **etcdadm**
![Version](https://img.shields.io/badge/version-5b496a72af3d80d64a16a650c85ce9a5882bc014-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiK0pzWGNJc01qaEVYTU9JcjY5MzdFTFVlSmV2aE1ESUVlODhKNHErSUNJSlkrV1o2bDlPS1hRU1BsWGJhNTZEVkNEYXVGeGRpRnJ4VkpjdFNiR2ZVQ21nPSIsIml2UGFyYW1ldGVyU3BlYyI6Ikh6dkhlYVh0QnE1TytCaU0iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[etcdadm](https://github.com/kubernetes-sigs/etcdadm) is a command-line tool for operating an etcd cluster. It downloads a specific etcd release, installs the binary, configures a systemd service, generates CA certificates, calls the etcd API to add (or remove) a member, and verifies that the new member is healthy. Its user experience is inspired by kubeadm.
