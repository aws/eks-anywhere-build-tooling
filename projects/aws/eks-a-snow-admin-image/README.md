# Snow EKS-A Admin image
This project builds an VM image compatible with Snow devices containing all necessary artifacts to manage EKS-A clusters.

It first uses `packer` to build an AMI and then it exports it to RAW format using `aws ec2 export-image`.

## Add more components
1. Add a bash script that installs your component under `providers`
1. Add a bash script that tests that your component has been installed under `provisioners/test` with the same name as installation script.
1. Reference your script in the build section of `snow-admin.pkr.hcl` using a `shell` provisioner.
	* If your script doesn't have any special requirements (like supporting reboots), use the same provisioner that references all the other install scripts.
	* If you need to reboot the machine from your script, use a new shell provisioner. Beware of order.
	* Don't add anything after the cleanup provisioner.

### Extra input
If your new component needs dynamic input:
1. Declare a `variable` in `snow-admin.pkr.hcl`.
1. Reference the variable with `${var.new-variable-name}`
1. If you need the variable in the installation script, pass it to the provisioner with and env var using `environment_vars`
1. In the `Makefile`, `export` a `PKR_VAR_new-variable-name` env var with the proper value in the `build-ami` target