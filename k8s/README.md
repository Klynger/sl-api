# Kubernetes manifests

What sl-api needs to run in a Kubernetes cluster. It lives here, next to the code, so everything about the app goes away with the app. The cluster itself (k3d config, image registry) lives in [personal-infra](https://github.com/Klynger/personal-infra), which only knows that sl-api exists.

For now these are applied by hand against the local k3d cluster from personal-infra.

## Namespace and Postgres

Everything goes in the `sl-api` namespace, including the app's own Postgres. The database credentials come from `k8s/postgres.env`, which is ignored by Git. Create it from the example, set a real password, and load it as a Secret:

```bash
cp k8s/postgres.env.example k8s/postgres.env   # then edit the password
kubectl apply -f k8s/namespace.yaml
kubectl -n sl-api create secret generic postgres --from-env-file=k8s/postgres.env
kubectl apply -R -f k8s/
kubectl -n sl-api rollout status statefulset/postgres
```

Inside the namespace the database is reachable as `postgres:5432`. Its data lives on the `data-postgres-0` volume claim, which is kept when the StatefulSet is deleted. Deleting the data is always a deliberate step: delete that claim, or the whole `sl-api` namespace, which takes everything in it.

Open a SQL shell:

```bash
kubectl -n sl-api exec -it postgres-0 -- psql -U shoppinglist_user -d shoppinglist_db
```
