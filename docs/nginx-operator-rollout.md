# Frontend NGINX compatibility preparation

Frontend Fdeployment manifests now set security.nginxCompatibility=true, privileged=false and runAsNonRoot=false. f-operator v0.2.2+ grants only CHOWN, SETUID, SETGID and NET_BIND_SERVICE to the standard root-master NGINX image, while keeping other capabilities dropped and privilege escalation/token mounting disabled. Backends are unchanged.

Do not merge this preparation branch into a deployment branch until the coordinated rollout: existing push workflows can build images and apply manifests automatically. Install the updated Fdeployment CRD, patch security on the live frontend CRs while preserving live versions/replicas/environment, then deploy the corrected operator. Merge this manifest change when the frontend deploy workflows should use the new option. Source manifests contain historical tags/counts and must not be applied blindly.

Full runbook and exact tested image digests: https://github.com/FabioKaelin/f-operator/blob/main/docs/nginx-compatibility.md . Existing development latest tags remain supported; version pinning is separate. This preparation does not change frontend images, application code or the running cluster.
