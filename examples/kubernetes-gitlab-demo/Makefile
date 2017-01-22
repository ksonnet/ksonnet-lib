gitlab-full.yml: $(wildcard gitlab/*.yml) $(wildcard gitlab-runner/*.yml)
		bash generate.bash

load-balancer-full.yml: $(wildcard load-balancer/*/*.yaml)
		bash generate.bash

apply: gitlab-full.yml load-balancer-full.yml
		kubectl apply -f load-balancer-full.yml
		kubectl apply -f gitlab-ns.yml
		kubectl apply -f gke/storage.yml
		kubectl apply -f gitlab-config.yml
		kubectl apply -f gitlab-full.yml
		kubectl apply -f ingress/gitlab-ingress.yml

delete: gitlab-full.yml
		kubectl delete -f gitlab-full.yml
		kubectl delete -f ingress/gitlab-ingress.yml

.PHONY: apply delete
