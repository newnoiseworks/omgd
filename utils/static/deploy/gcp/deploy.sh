#!/bin/bash

export CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE=~/.config/gcloud/application_default_credentials.json

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "truncate -s 0 /var/log/docker.log" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

# scp'ing the docker-compose file in case it's the first deploy and we don't have one, so the subsequent down command won't error out
gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --force-key-file-overwrite ../../../servers/docker-compose.yml ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}:

# $GCP_UPDATE_REMOVE_VOLUME not set, then...
if [[ -z "${GCP_UPDATE_REMOVE_VOLUME}" ]]; then
    gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "\$PWD:\$PWD" -w="\$PWD" docker/compose:1.27.0 down" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}
else
    gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "\$PWD:\$PWD" -w="\$PWD" docker/compose:1.27.0 down -v" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}
fi

gcloud compute scp --project ${GCP_PROJECT} --zone ${GCP_ZONE} --recurse --force-key-file-overwrite ../../../servers/* ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}:

gcloud compute ssh --project ${GCP_PROJECT} --zone ${GCP_ZONE} --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "\$PWD:\$PWD" -w="\$PWD" docker/compose:1.27.0 up -d ${OMGD_SERVER_SERVICES}" ${OMGD_PROJECT}-omgd-dev-instance-${OMGD_PROFILE}

