# OpenEBS package miscellaneous files

## Scripts
- **demo.sh** - _Runs through the full presentation demo. Installs openEBS, configures cStor, creates pvc, deploys test application, tests resizing and data retention.
- **helminst.sh** - shorthand to install the helm chart
- **install-and-test.sh** - Installs the helm chart, configures cStor, sets up a pvc, and runs a basic test pod.

## Configs
- **eksa-seifsall-cluster.yaml** - Working cluster config file for a BareMetal cluster with 1 CP node and 2 worker nodes. IP addresses may need to be slightly modified
- **configs/cstor-configs** - Config files for setting up cStor
  - **cspc.yaml** - Sets up a cStor pool cluster with 2 nodes. This will need to be modified for each cluster it is installed on
  - **cstor-csi-disk.yaml** - Sets up a storage class using the cStor pool cluster created by _cspc.yaml_
- **configs/pvc** - Config files for setting up PVCs
  - **pvc.yaml** - Creates a 100M PVC called small-pvc using the cStor storage class
- **configs/sample** - Config files for a sample busybox application
  - **busybox.yaml** - Creates a simple pod that writes "Hello from <node name>" with a timestamp to a file called _greet.txt_. This will need to be modified to have the correct node name for each cluster.
- **configs/snapshot** - Config files for setting up snapshotter
  - **snapshot_class.yaml** - Creates a storage class for snapshotter
  - **snapshot.yaml** - Creates a pvc snapshot of a pvc named "small-pvc".