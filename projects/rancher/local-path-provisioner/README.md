## **Local Path Provisioner**
![Version](https://img.shields.io/badge/version-v0.0.32-blue)
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiNmJlc3diN0NwYzhjZUtZc2wvSVdQVk16aFJTNmtvSnpPWmZhMjZaM0tkNU5QZCtsQXluamlQWVd6cVJNTVRjTmM2ZVAzUnlFTVozUVA4Um5XZTJpNXlrPSIsIml2UGFyYW1ldGVyU3BlYyI6Iktaam5IZ3JCVFBheXMydDIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

[Local Path Provisioner](https://github.com/rancher/local-path-provisioner) provides a way for the Kubernetes users to utilize the local storage in each node. Based on the user configuration, the Local Path Provisioner will create `hostPath`-based persistent volume on the node automatically. It utilizes the features introduced by Kubernetes Local Persistent Volume feature, but makes it a simpler solution than the built-in `local` volume feature in Kubernetes. Its advantage over the Local PV feature is its ability to perform dynamic provisioning of volumes using `hostPath`.

You can find the latest version of this image [on ECR Public Gallery](https://gallery.ecr.aws/eks-anywhere/rancher/local-path-provisioner).