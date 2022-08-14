eksctl create cluster \
    --name test-cluster \
    --region us-west-2 \
    --nodegroup-name linux-nodes \
    --node-type t2.micro \
    --nodes 2