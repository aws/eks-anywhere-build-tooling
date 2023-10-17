# EKS-A Admin image
![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiSGtiejcreFBSRjBuVzV1cTlIb1BTa09KaGt2Z3lOZ1JLblBLQ0hwU1phQ0JyUjRKcmZ5SlB6UThXZDJhZ3JhZ3U1cFVlZ1BDNFJvS1FwcjlUMUtWRXh3PSIsIml2UGFyYW1ldGVyU3BlYyI6InBBTFpDN2xRRGNyMFZIRXIiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

This project uses Packer to build EKS-A Admin VM images containing all necessary artifacts to create and manage EKS-A clusters.

## Add more components
1. Add a bash script that installs your component under `providers`
1. Add a bash script that tests that your component has been installed under `provisioners/test` with the same name as installation script.
1. Reference your script in the build section of `eks-a-admin.pkr.hcl` using a `shell` provisioner.
	* If your script doesn't have any special requirements (like supporting reboots), use the same provisioner that references all the other install scripts.
	* If you need to reboot the machine from your script, use a new shell provisioner. Beware of order.
	* Don't add anything after the cleanup provisioner.
	* To enable a provisioner for only your image add `only = ["vsphere-iso.ubuntu"]` to your provisioner block.

### Extra input
If your new component needs dynamic input:
1. Declare a `variable` in `eks-a-admin.pkr.hcl`.
1. Reference the variable with `${var.new-variable-name}`
1. If you need the variable in the installation script, pass it to the provisioner with and env var using `environment_vars`
1. In the `Makefile`, `export` a `PKR_VAR_new-variable-name` env var with the proper value in the `build-ami` target

### Admin Images

1. **Snow Device Admin AMI Image (source.amazon-ebs.ubuntu):** Builds an EKS-A Admin (ami) image compatible with Snow devices. It first uses `packer` to build an AMI and then it exports it to RAW format using `aws ec2 export-image`.
1. **Tinkerbell CI Admin OVA Image (source.vsphere-iso.ubuntu):** Builds an EKS-A Admin (ova) image which serves as admin machine during Tinkerbell e2e tests. Build's an OVA using `packer` and then exports the OVA to s3.