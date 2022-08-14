aws ecr create-repository \
    --repository-name e-commerce-api \
    --image-scanning-configuration scanOnPush=true \
    --region us-west-2