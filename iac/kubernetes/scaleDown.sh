gcloud container clusters resize cluster --node-pool cadence --num-nodes 0
gcloud container clusters resize cluster --node-pool cadence-web --num-nodes 0
gcloud container clusters resize cluster --node-pool cadence-frontend --num-nodes 0
gcloud container clusters resize cluster --node-pool cadence-worker --num-nodes 0
gcloud container clusters resize cluster --node-pool cadence-matching --num-nodes 0
gcloud container clusters resize cluster --node-pool cadence-history --num-nodes 0
gcloud container clusters resize cluster --node-pool cadence-cassandra --num-nodes 0
gcloud container clusters resize cluster --node-pool default-pool --num-nodes 0