gcloud beta container --project "e-commerce-332709" clusters create "cluster" --zone "us-east1-b" --no-enable-basic-auth /
    --cluster-version "1.21.5-gke.1302" --release-channel "regular" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
     --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
     "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write",/
     "https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol",/
     "https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" /
     --max-pods-per-node "110" --num-nodes "3" --logging=SYSTEM,WORKLOAD --monitoring=SYSTEM --enable-ip-alias /
     --network "projects/e-commerce-332709/global/networks/default" --subnetwork "projects/e-commerce-332709/regions/us-east1/subnetworks/default" /
     --no-enable-intra-node-visibility --default-max-pods-per-node "110" --no-enable-master-authorized-networks /
     --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver --enable-autoupgrade --enable-autorepair /
     --max-surge-upgrade 1 --max-unavailable-upgrade 0 --enable-shielded-nodes --node-locations "us-east1-b" && /

     gcloud beta container --project "e-commerce-332709" node-pools create "cadence" --cluster "cluster" --zone "us-east1-b" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
     --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
     "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring",/
     "https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" /
     --num-nodes "1" --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --max-pods-per-node "110" --node-locations "us-east1-b" && /

     gcloud beta container --project "e-commerce-332709" node-pools create "cadence-frontend" --cluster "cluster" --zone "us-east1-b" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
     --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
     "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring",/
     "https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" /
     --num-nodes "1" --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --max-pods-per-node "110" --node-locations "us-east1-b" && /

     gcloud beta container --project "e-commerce-332709" node-pools create "cadence-cassandra" --cluster "cluster" --zone "us-east1-b" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
     --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
     "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring",/
     "https://www.googleapis.com/auth/servicecontrol", "https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" /
     --num-nodes "1" --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --max-pods-per-node "110" --node-locations "us-east1-b" && /

     gcloud beta container --project "e-commerce-332709" node-pools create "cadence-worker" --cluster "cluster" --zone "us-east1-b" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
     --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
     "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring", /
     "https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" /
      --num-nodes "1" --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --max-pods-per-node "110" --node-locations "us-east1-b" && /

      gcloud beta container --project "e-commerce-332709" node-pools create "cadence-web" --cluster "cluster" --zone "us-east1-b" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
      --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
      "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring",/
      "https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append"/
       --num-nodes "1" --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --max-pods-per-node "110" --node-locations "us-east1-b" && /

      gcloud beta container --project "e-commerce-332709" node-pools create "cadence-history" --cluster "cluster" --zone "us-east1-b" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
      --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
      "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring",/
      "https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" /
      --num-nodes "1" --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --max-pods-per-node "110" --node-locations "us-east1-b" && /

      gcloud beta container --project "e-commerce-332709" node-pools create "cadence-matching" --cluster "cluster" --zone "us-east1-b" --machine-type "e2-medium" --image-type "COS_CONTAINERD" /
      --disk-type "pd-standard" --disk-size "100" --metadata disable-legacy-endpoints=true --scopes /
      "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring",/
      "https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append"/
       --num-nodes "1" --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --max-pods-per-node "110" --node-locations "us-east1-b"
