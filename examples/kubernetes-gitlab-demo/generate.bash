#!/bin/bash

set -e

scope_files() {
  for file in "$@"; do
    cat "$file"
    echo "---"
  done
}

scope_files \
  gitlab/redis-* \
  gitlab/postgresql-* \
  gitlab/gitlab-* \
  gitlab-runner/gitlab-runner-* \
  > gitlab-full.yml

scope_files \
  load-balancer/lego/* \
  load-balancer/nginx/* \
  > load-balancer-full.yml

scope_files ingress/* \
  > gitlab-ingress.yml

scope_files \
  load-balancer/lego/* \
  load-balancer/nginx/* \
  gitlab-ns.yml \
  gitlab-config.yml \
  gke/storage.yml \
  gitlab/redis-* \
  gitlab/postgresql-* \
  gitlab/gitlab-* \
  gitlab-runner/gitlab-runner-* \
  ingress/* \
  > gitlab-all.yml


if [[ -z "$GITLAB_GKE_IP" || -z "$GITLAB_GKE_DOMAIN" || -z "$GITLAB_LEGO_EMAIL" ]]; then
    echo "A needed variable is not exported. Check that GITLAB_GKE_IP, GITLAB_GKE_DOMAIN, and GITLAB_LEGO_EMAIL environment variables are set. Cannot replace & rename."
    exit 1
fi

dashed_domain="${GITLAB_GKE_DOMAIN//./-}"
echo "Using gitlab-${dashed_domain}.yml"
cp gitlab-all.yml gitlab-${dashed_domain}.yml

sed -i.bak "
    s/@GITLAB_GKE_IP@/${GITLAB_GKE_IP}/g
    s/@GITLAB_GKE_DOMAIN@/${GITLAB_GKE_DOMAIN}/g
    s/@GITLAB_LEGO_EMAIL@/${GITLAB_LEGO_EMAIL}/g
" gitlab-${dashed_domain}.yml
rm gitlab-${dashed_domain}.yml.bak
