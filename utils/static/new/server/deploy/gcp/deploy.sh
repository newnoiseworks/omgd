#!/bin/bash

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "truncate -s 0 /var/log/docker.log" nakama-instance

# scp'ing the docker-compose file in case it's the first deploy and we don't have one, so the subsequent down command won't error out
gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --force-key-file-overwrite ../../docker-compose.yml nakama-instance:

# $GCP_UPDATE_REMOVE_VOLUME not set, then...
if [[ -z "${GCP_UPDATE_REMOVE_VOLUME}" ]]; then
    gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "\$PWD:\$PWD" -w="\$PWD" docker/compose:1.24.0 down" nakama-instance
else
    gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "\$PWD:\$PWD" -w="\$PWD" docker/compose:1.24.0 down -v" nakama-instance
fi

gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --recurse --force-key-file-overwrite ../../nakama nakama-instance:

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "\$PWD:\$PWD" -w="\$PWD" docker/compose:1.24.0 up -d" nakama-instance

