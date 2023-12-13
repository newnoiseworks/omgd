#!/bin/bash

export CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=~/.config/gcloud/application_default_credentials.json

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "sudo curl -sSL https://github.com/docker/compose/releases/download/v2.23.3/docker-compose-linux-x86_64 -o /var/lib/google/docker-compose && sudo chmod o+x /var/lib/google/docker-compose && mkdir -p ~/.docker/cli-plugins && sudo ln -sf /var/lib/google/docker-compose ~/.docker/cli-plugins/docker-compose" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "truncate -s 0 /var/log/docker.log" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

# scp'ing the docker-compose file in case it's the first deploy and we don't have one, so the subsequent down command won't error out
gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --force-key-file-overwrite ../../../servers/docker-compose.yml ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}:

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker compose down" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --recurse --force-key-file-overwrite ../../../servers/* ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}:

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker compose up -d ${OMGD_SERVER_SERVICES}" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

