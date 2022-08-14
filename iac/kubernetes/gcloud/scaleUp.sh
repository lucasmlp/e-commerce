gcloud container clusters resize cluster --node-pool default-pool --num-nodes 1
gcloud container clusters resize cluster --node-pool cadence-cassandra --num-nodes 1
gcloud container clusters resize cluster --node-pool cadence-history --num-nodes 1
gcloud container clusters resize cluster --node-pool cadence-matching --num-nodes 1
gcloud container clusters resize cluster --node-pool cadence-worker --num-nodes 1
gcloud container clusters resize cluster --node-pool cadence-frontend --num-nodes 1
gcloud container clusters resize cluster --node-pool cadence-web --num-nodes 1
gcloud container clusters resize cluster --node-pool cadence --num-nodes 1