# tier 1
# GCP_FUNCTION_ZONES=('us-west1' 'us-central1' 'us-east1' 'us-east4' 'europe-west1' 'europe-west2' 'asia-east1' 'asia-east2' 'asia-northeast1' 'asia-northeast2')
# tier 2
# GCP_FUNCTION_ZONES=('us-west2' 'us-west3' 'us-west4' 'northamerica-northeast1' 'southamerica-east1' 'europe-west3' 'europe-west6' 'europe-central2' 'australia-southeast1' 'asia-south1' 'asia-southeast1' 'asia-southeast2' 'asia-northeast3')
GCP_FUNCTION_ZONES=('asia-northeast3')

for zone in "${GCP_FUNCTION_ZONES[@]}"
do
  gcloud functions deploy HandleProxyRequest --runtime go116 --trigger-http --allow-unauthenticated --region $zone --memory 128MB
done