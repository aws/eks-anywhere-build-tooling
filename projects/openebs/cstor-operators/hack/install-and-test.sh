printf "OpenEBS Basic Test Script"
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
printf "\nkubectl apply -f configs/cstor_configs/cspc.yaml\n"
kubectl apply -f configs/cstor_configs/cspc.yaml
printf "\nWaiting for CSPC to be ready... (30s)\n"
sleep 30

printf "\nCreate cStor storage class\n"
printf "\nkubectl apply -f configs/cstor_configs/cstor-csi-disk.yaml\n"
kubectl apply -f configs/cstor_configs/cstor-csi-disk.yaml
printf "\n waiting for storage class to be ready (10s)\n"
sleep 10

printf "\ncStor setup complete!\n"

printf "\n\nStarting PVC setup\n"
printf "\nCreating PVC\n"
printf "\nkubectl apply -f configs/pvc/pvc.yaml\n"
kubectl apply -f configs/pvc/pvc.yaml
printf "\nWaiting for PVC to be ready.... (10s)\n"
sleep 10
printf "\nPVC setup complete!\n"


printf "\n\nTest application pod setup\n"
sleep 15
printf "\nkubectl apply -f configs/sample/busybox2.yaml\n"
kubectl apply -f configs/sample/busybox2.yaml
printf "\nwaiting for pod to be ready (30s)\n"
sleep 30
printf "\nRun test exec to print out current contents of greet.txt\n"
printf "\nkubectl exec -it busybox2 -- cat /mnt/openebs-csi/greet.txt\n"
kubectl exec -it busybox2 -- cat /mnt/openebs-csi/greet.txt
printf "\nTest pod deployment complete!\n"