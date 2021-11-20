#!/bin/bash

# Exit on any error
set -e

# Config google-cloud-sdk and kubectl
sudo apt update
sudo apt install lsb-release
export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)"
echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
sudo apt-get update && sudo apt-get install google-cloud-sdk kubectl

# Setup google-cloud-sdk auth
# echo $GCLOUD_SERVICE_KEY | base64 --decode -i > ${HOME}/account-auth.json
# export GOOGLE_APPLICATION_CREDENTIALS=${HOME}/account-auth.json

# gcloud auth activate-service-account --key-file ${HOME}/account-auth.json

# Setup project and cluster
gcloud config set project $PROJECT_NAME
gcloud --quiet config set container/cluster $CLUSTER_NAME

# Reading the zone from the env var is not working so we set it here
# gcloud config set compute/zone ${CLOUDSDK_COMPUTE_ZONE}
gcloud --quiet container clusters get-credentials $CLUSTER_NAME

# Building
docker build -t gcr.io/${PROJECT_NAME}/journey-service-workflow-worker:$CIRCLE_SHA1 --target workflowworker .

# Tagging
docker tag gcr.io/${PROJECT_NAME}/journey-service-workflow-worker:$CIRCLE_SHA1 gcr.io/${PROJECT_NAME}/journey-service-workflow-worker:latest

# Config docker to use google repository
gcloud auth configure-docker --quiet

# Push images
docker push gcr.io/${PROJECT_NAME}/journey-service-workflow-worker

# Permissions
sudo chown -R circleci:circleci /home/circleci/.config
sudo chown -R circleci:circleci /home/circleci/.kube

# Deploy
kubectl patch deployment journey-service-workflow-worker -p '{"spec":{"template":{"spec":{"containers":[{"name":"journey-service-workflow-worker","image":"gcr.io/'"$PROJECT_NAME"'/journey-service-workflow-worker:'"$CIRCLE_SHA1"'"}]}}}}'