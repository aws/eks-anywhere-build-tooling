printf "OpenEBS Demo Script"
printf "\n\nInstalling OpenEBS with custom helm chart"
printf "\nhelm install openebs-cstor oci://public.ecr.aws/g4p7p1c7/openebs/cstor --namespace openebs-cstor --version 3.3.0-75eac31e866fec47b365f8601fae64b6ee1e636a --set sourceRegistry=public.ecr.aws/g4p7p1c7 --create-namespace\n"
helm install openebs-cstor oci://public.ecr.aws/g4p7p1c7/openebs/cstor --namespace openebs-cstor --version 3.3.0-75eac31e866fec47b365f8601fae64b6ee1e636a --set sourceRegistry=public.ecr.aws/g4p7p1c7 --create-namespace

printf "\nWaiting for pods to be ready.... (50s)\n"
sleep 50
printf "\nkubectl get pods -n openebs-cstor\n"
kubectl get pods -n openebs-cstor


printf "\n\nChecking if BDs were created\n"
printf "\nkubectl get bd -n openebs-cstor\n"
kubectl get bd -n openebs-cstor
printf "\nInitial setup complete!\n"

printf "\n\nStarting cStor setup\n"
printf "\nCreating CSPC\n"
printf "\nShowing CSPC config\n"
printf "\ncat configs/cstor_configs/cspc.yaml\n"
cat configs/cstor_configs/cspc.yaml
sleep 15
printf "\nkubectl apply -f configs/cstor_configs/cspc.yaml\n"
kubectl apply -f configs/cstor_configs/cspc.yaml
printf "\nWaiting for CSPC to be ready... (30s)\n"
sleep 30

printf "\nkubectl get cspc -n openebs-cstor\n"
kubectl get cspc -n openebs-cstor

printf "\nCreate cStor storage class\n"
printf "\nShowing Storage Class config\n"
printf "\ncat configs/cstor_configs/cstor-csi-disk.yaml\n"
cat configs/cstor_configs/cstor-csi-disk.yaml
sleep 15
printf "\nkubectl apply -f configs/cstor_configs/cstor-csi-disk.yaml\n"
kubectl apply -f configs/cstor_configs/cstor-csi-disk.yaml
printf "\n waiting for storage class to be ready (10s)\n"
sleep 10
printf "\nkubectl get sc\n"
kubectl get sc

printf "\ncStor setup complete!\n"


printf "\n\nStarting PVC setup\n"
printf "\nCreating PVC\n"
printf "\nShowing PVC config\n"
printf "\ncat configs/pvc/pvc.yaml\n"
cat configs/pvc/pvc.yaml
sleep 15
printf "\nkubectl apply -f configs/pvc/pvc.yaml\n"
kubectl apply -f configs/pvc/pvc.yaml
printf "\nWaiting for PVC to be ready.... (10s)\n"
sleep 10

printf "\nkubectl get pvc\n"
kubectl get pvc
printf "\nPVC setup complete!\n"


printf "\n\nTest application pod setup\n"
printf "\nShowing pod config\n"
printf "\ncat configs/sample/busybox2.yaml\n"
cat configs/sample/busybox2.yaml
sleep 15
printf "\nkubectl apply -f configs/sample/busybox2.yaml\n"
kubectl apply -f configs/sample/busybox2.yaml
printf "\nwaiting for pod to be ready (30s)\n"
sleep 30
printf "\nRun test exec to print out current contents of greet.txt\n"
printf "\nkubectl exec -it busybox2 -- cat /mnt/openebs-csi/greet.txt\n"
kubectl exec -it busybox2 -- cat /mnt/openebs-csi/greet.txt
printf "\nInitial test pod deployment complete!\n"
sleep 20
printf "\n\nPVC Resizing test\n"
printf "\nDeploying 90M of data to 100M PVC\n"
printf "\nkubectl exec -it busybox2 -- sh -c 'dd if=/dev/urandom bs=1M count=90 | tee /mnt/openebs-csi/test.data | sha256sum'\n"
kubectl exec -it busybox2 -- sh -c 'dd if=/dev/urandom bs=1M count=90 | tee /mnt/openebs-csi/test.data | sha256sum'

printf "\nEdit PVC to increase size from 100M to 150M\n"
printf "\nkubectl edit pvc small-pvc\n"
kubectl edit pvc small-pvc
printf "\n Wait for resize to complete.... (60s)\n"
sleep 60
printf "\nkubectl get pvc\n"
kubectl get pvc
printf "\nWrite 50M more data (total 140M), exceeding initial 100M capacity.\n"
printf "\nOriginal data sha256 sum: \n"
kubectl exec -it busybox2 -- sha256sum /mnt/openebs-csi/test.data
printf "\nWriting new data\n"
printf "\nkubectl exec -it busybox2 -- sh -c 'dd if=/dev/urandom bs=1M count=50 | tee /mnt/openebs-csi/test2.data | sha256sum'\n"
kubectl exec -it busybox2 -- sh -c 'dd if=/dev/urandom bs=1M count=50 | tee /mnt/openebs-csi/test2.data | sha256sum'
printf "\nCheck previous data sha256sum to make sure it is still intact after exceeding initial storage limit\n"
printf "\nkubectl exec -it busybox2 -- sha256sum /mnt/openebs-csi/test.data\n"
kubectl exec -it busybox2 -- sha256sum /mnt/openebs-csi/test.data
printf "\nResizing test complete!\n"

printf "\n\nPod deletion data retention test\n"
printf "\nDelete busybox pod\n"
printf "\nkubectl delete pod busybox2\n"
kubectl delete pod busybox2
printf "\nSchedule busybox2 to the other worker node\n"
printf "\nnano configs/sample/busybox2.yaml\n"
nano configs/sample/busybox2.yaml
printf "\nDeploy newly edited pod\n"
printf "\nkubectl apply -f configs/sample/busybox2.yaml\n"
kubectl apply -f configs/sample/busybox2.yaml
printf "\nWait for pod to be ready.... (30s)\n"
sleep 30
printf "\nCheck if pod is healthy\n"
printf "\nkubectl get pods\n"
kubectl get pods
printf "\nChecking if pod prints out retained data from PVC after being scheduled to another node\n"
printf "\nkubectl exec -it busybox2 -- cat /mnt/openebs-csi/greet.txt\n"
kubectl exec -it busybox2 -- cat /mnt/openebs-csi/greet.txt
