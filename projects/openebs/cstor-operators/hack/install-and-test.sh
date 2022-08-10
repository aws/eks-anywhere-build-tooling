#!/bin/bash

# uninstall lingering resources from openebs
echo ">>>>>>>>>>>>>>> DELETE RESOURCES <<<<<<<<<<<<<<<"
kubectl delete pod busybox
kubectl delete pvc cstor-pvc
kubectl delete sc cstor-csi-disk
kubectl delete cspc cstor-disk-pool -n openebs
sleep 100
echo "Done!!!"
# Start by uninstalling openebs from the cluster, if it exists
echo "\n>>>>>>>>>>>>>>> RUN OPENEBS UNINSTALL SCRIPT <<<<<<<<<<<<<<<"
bash uninstall.sh
sleep 120
echo "Done!!!"

# Install helm chart with cstor operator using custom values.yaml
echo "\n>>>>>>>>>>>>>>> DELETE HELM REPO <<<<<<<<<<<<<<<"
helm repo remove openebs
echo ">>>>>>>>>>>>>>> ADD HELM REPO <<<<<<<<<<<<<<<"
helm repo add openebs https://openebs.github.io/charts
echo ">>>>>>>>>>>>>>> UPDATE HELM REPO <<<<<<<<<<<<<<<"
helm repo update
echo ">>>>>>>>>>>>>>> INSTALL HELM CHART <<<<<<<<<<<<<<<"
helm install openebs --namespace openebs openebs/openebs --set cstor.enabled=true --create-namespace -f values.yaml --wait=true
echo "Helm install done!!!"
sleep 60

# Set up cStor storage engine
echo "\n>>>>>>>>>>>>>>> Creating cStor Storage Pools <<<<<<<<<<<<<<<"
kubectl apply -f ./seifsall/cstor_configs/cspc.yaml
echo ">>>>>>>>>>>>>>> Creating cStor Storage Classes <<<<<<<<<<<<<<<"
kubectl apply -f ./seifsall/cstor_configs/cstor-csi-disk.yaml
sleep 60
echo "Done!!!"

# Deploy busybox application
echo "\n>>>>>>>>>>>>>>> Creating PVC <<<<<<<<<<<<<<<"
kubectl apply -f ./seifsall/pvc/pvc.yaml
echo ">>>>>>>>>>>>>>> Deploying test app pod <<<<<<<<<<<<<<<"
kubectl apply -f ./seifsall/sample/busybox.yaml --wait=true
echo "App deployed!!!"
echo "Waiting for app to be ready...."
sleep 200
# Run test app and check output
echo "-----------------------"
echo "-----------------------"
echo "TEST APPLICATION OUTPUT:"
kubectl exec -it busybox -- cat /mnt/openebs-csi/date.txt
