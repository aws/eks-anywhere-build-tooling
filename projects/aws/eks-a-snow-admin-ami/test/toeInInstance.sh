
set -e
set -o pipefail

AMI_ID=$1
INSTANCE_TYPE=$2 # t2.large
KEY_NAME=$3 # MyKeyPair
KEY_PATH=$4 # ~/.ssh/mykey.pem
USER=$5 # user to ssh into instance, probably dependeing on OS
DOCUMENTS_ARG=$6 # documents to pass to awstoe
PHASES=$7 # phases to run with awstoe
PARAMETERS=$8 # parameters to pass to awstoe run (optional)

echo "Starting test instance"
instance_id=$(aws ec2 run-instances --image-id $AMI_ID --count 1 --instance-type $INSTANCE_TYPE --key-name $KEY_NAME --output text --query 'Instances[*].InstanceId')
echo "Waiting for $instance_id to be running"
aws ec2 wait instance-running --instance-ids $instance_id

public_dns=$(aws ec2 describe-instances --instance-ids $instance_id | jq -r '.Reservations[0].Instances[0].PublicDnsName' -)
user_host=$USER@$public_dns

echo "Waiting for ssh in $instance_id to be ready"
# TODO: make this more robust with some kind of retry and timeout
ssh -i $KEY_PATH -o StrictHostKeyChecking=no $user_host "ls 2>&1"

echo "Instance $instance_id is ready"

#make sure awstoe is installed
mkdir -p bin
[[ ! -f bin/awstoe ]] && echo "Downlaoding awstoe" && curl -L --silent https://awstoe-us-east-1.s3.us-east-1.amazonaws.com/latest/linux/amd64/awstoe > ./bin/awstoe
chmod +x bin/awstoe

DOCUMENTS=(${DOCUMENTS_ARG//,/ })
WORKING_DIR=/usr/src/project
LOGS_DIR=$WORKING_DIR/logs
mkdir -p logs

echo "Copying binaries"
scp -q -i $KEY_PATH -o StrictHostKeyChecking=no -r bin $user_host:bin

echo "Copying documents"
for i in "${DOCUMENTS[@]}"
do
	ssh -i $KEY_PATH -o StrictHostKeyChecking=no $user_host "mkdir -p $(dirname $i)"
	scp -q -i $KEY_PATH -o StrictHostKeyChecking=no $i $user_host:$i
done

awstoe_command="./bin/awstoe run --documents $DOCUMENTS_ARG --phases $PHASES"
[[ ! -z "$PARAMETERS" ]] && awstoe_command="$awstoe_command --parameters $PARAMETERS"

echo "Running awstoe in instance"
awstoe_response=$(ssh -i $KEY_PATH -o StrictHostKeyChecking=no $user_host "$awstoe_command")

awstoe_status=$(echo $awstoe_response | jq -r '.status' -)
remote_log_path=$(echo $awstoe_response | jq -r '.logUrl' -)
logs_folder=$(basename $remote_log_path)

mkdir -p logs

local_logs_path=logs/$logs_folder

echo "Copying awstoe logs to local disk"
scp -q -i $KEY_PATH -o StrictHostKeyChecking=no -r $user_host:$remote_log_path $local_logs_path

echo "You can find awstoe logs at $local_logs_path"

if [[ "$awstoe_status" == "success" ]]; then
	echo "Awstoe run finished succesfully"
	echo "Terminating instance $instance_id"
	aws ec2 terminate-instances --instance-ids $instance_id
	return_code=0
else 
	echo "awstoe failed with status $awstoe_status"
	echo "Instance $instance_id won't be terminated to allow debugging"
	echo "To ssh into it: ssh -i $KEY_PATH -o StrictHostKeyChecking=no $user_host"
	echo "To terminate the instance: aws ec2 terminate-instances --instance-ids $instance_id"
	echo "awstoe logs:"
	cat $local_logs_path/console.log
	return_code=1
fi

exit $return_code
