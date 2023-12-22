#!/bin/bash

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "sudo curl -sSL https://github.com/docker/compose/releases/download/v2.23.3/docker-compose-linux-x86_64 -o /var/lib/google/docker-compose && sudo chmod o+x /var/lib/google/docker-compose && mkdir -p ~/.docker/cli-plugins && sudo ln -sf /var/lib/google/docker-compose ~/.docker/cli-plugins/docker-compose" omgd-sa@${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "truncate -s 0 /var/log/docker.log" omgd-sa@${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --force-key-file-overwrite ../../../servers/docker-compose.yml omgd-sa@${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}:

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker compose down" omgd-sa@${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --recurse --force-key-file-overwrite ../../../servers/* omgd-sa@${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}:

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker compose up -d ${OMGD_SERVER_SERVICES}" omgd-sa@${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

