# Kubernetes manifests

What sl-api needs to run in a Kubernetes cluster. It lives here, next to the code, so everything about the app goes away with the app. The cluster itself (k3d config, image registry) lives in [personal-infra](https://github.com/Klynger/personal-infra), which only knows that sl-api exists.

For now these are applied by hand against the local k3d cluster from personal-infra, always with server-side apply (`--server-side`). The API server then tracks which fields kubectl owns and compares only those, so fields it fills in itself, like the defaults inside the Postgres volume claim template, no longer show up as changes on every apply. Server-side apply prints `serverside-applied` whether or not anything changed: `kubectl diff --server-side -R -f k8s/` shows what an apply would actually change.

## Config and secrets

Plain settings live in the `sl-api` ConfigMap (`k8s/config.yaml`), committed like any other manifest. Each secret value is its own Secret, named `production--sl-api--<name>` with the value under the key `value`, following the same naming as a cloud secret manager so the move to one later is one-to-one. Manifests read them one at a time:

```yaml
valueFrom:
  secretKeyRef:
    name: production--sl-api--db-password
    key: value
```

Secret values are not committed. They live in `k8s/secrets.env`, ignored by Git, one `<secret-name>=<value>` per line. Create it from the example and set real values, then load every line as its own Secret:

```bash
cp k8s/secrets.env.example k8s/secrets.env   # then edit the values
kubectl apply --server-side -f k8s/namespace.yaml
grep -v '^#' k8s/secrets.env | while IFS='=' read -r name value; do
  [ -n "$name" ] && kubectl -n sl-api create secret generic "$name" --from-literal=value="$value" \
    --dry-run=client -o yaml | kubectl apply --server-side -f -
done
```

Running the loop again updates the Secrets in place.

## Postgres

The app has its own Postgres in the `sl-api` namespace. The API and migrations connect with the same user Postgres creates on first start. Postgres reads its user, password and database name only when it first creates its data folder, so changing those values later changes the Secrets but not the database.

```bash
kubectl apply --server-side -R -f k8s/
kubectl -n sl-api rollout status statefulset/postgres
```

Inside the namespace the database is reachable as `postgres:5432`. Its data lives on the `data-postgres-0` volume claim, which is kept when the StatefulSet is deleted. Deleting the data is always a deliberate step: delete that claim, or the whole `sl-api` namespace, which takes everything in it.

Open a SQL shell:

```bash
kubectl -n sl-api exec -it postgres-0 -- psql -U shoppinglist_user -d shoppinglist_db
```

## API and migrations

The API runs as the `sl-api` Deployment, reachable inside the cluster as the `sl-api` Service on port 80. From this machine it is reachable at `http://sl-api.localhost:8080`: port 8080 goes to Traefik (the k3d port mapping in personal-infra), and the `sl-api` Ingress routes that hostname to the Service. Any `*.localhost` name resolves to this machine, so each app in the cluster can have its own. Migrations run as a Job named `migrate-<tag>`, separate from the API, which does not wait for them: every migration must keep working with the code already running (expand/contract).

Both use the image `registry.localhost:5000/sl-api:<tag>`, where the tag is the commit the image was built from. Build that exact commit, not the working tree, so the tag always matches the code inside:

```bash
TAG=$(git rev-parse --short HEAD)
git archive --format=tar "$TAG" | docker build -t "localhost:5000/sl-api:$TAG" -
docker push "localhost:5000/sl-api:$TAG"
```

### Deploying a new version

1. Build and push the new commit as above.
2. Set the new tag in `k8s/api/deployment.yaml` (the image) and in `k8s/migrate/job.yaml` (the image and the Job name).
3. `kubectl diff --server-side -R -f k8s/`, then `kubectl apply --server-side -R -f k8s/`.

The new Job runs `migrate up`, which applies only migrations the database does not have yet. The Deployment replaces the API Pod once the new one passes its `/livez` check. Old Jobs stay until deleted (`kubectl -n sl-api delete job migrate-<old-tag>`).

### Rolling back

Set the previous tag back, through a revert of the bump commit, and apply. The API goes back to the old image. The schema stays as it is, which is safe because migrations only add what older code ignores. Undoing a migration is a separate, deliberate step, and needs the newer image, since only it contains that migration's `down`.
