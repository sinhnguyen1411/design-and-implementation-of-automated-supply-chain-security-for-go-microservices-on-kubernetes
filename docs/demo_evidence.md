# Demo Evidence - Local Signed Demo on Docker Desktop Kubernetes

## Environment
- Kubernetes context: `docker-desktop`
- Kubernetes namespace: `stock-trading`
- Cluster nodes:
  - `desktop-control-plane`
  - `desktop-worker`
- Demo automation entrypoint: `scripts/local_signed_demo.ps1`

## Image Digest
Image:
`ttl.sh/stock-trading-d4ec05c5f73c@sha256:42e4b1f9ff6e20ec9ec3885614552472b680954f818456caa7aad6b675d36724`

SBOM digest:
`7C1EB2358B426D9CC031B954A7CC402770A5DA0DB6CB83B507967055AB826A38`

## Deployed Workload Status
```text
NAME                           READY   UP-TO-DATE   AVAILABLE   AGE    CONTAINERS     IMAGES                                                                                                      SELECTOR
deployment.apps/user-service   1/1     1            1           4m9s   user-service   ttl.sh/stock-trading-d4ec05c5f73c@sha256:42e4b1f9ff6e20ec9ec3885614552472b680954f818456caa7aad6b675d36724   app.kubernetes.io/name=user-service

NAME                               READY   STATUS    RESTARTS   AGE    IP            NODE             NOMINATED NODE   READINESS GATES
pod/user-service-fddcb8687-dzjrh   1/1     Running   0          4m7s   10.244.1.12   desktop-worker   <none>           <none>

NAME                   TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)               AGE     SELECTOR
service/user-service   ClusterIP   10.96.158.207   <none>        18080/TCP,19090/TCP   4m11s   app.kubernetes.io/name=user-service
```

## Pod Verification Evidence
Key annotations observed on the running pod:
- `kyverno.io/verify-images: {"ttl.sh/stock-trading-d4ec05c5f73c@sha256:42e4b1f9ff6e20ec9ec3885614552472b680954f818456caa7aad6b675d36724":"pass"}`
- `security.grype.io/high_critical: 0`
- `security.stock-trading.dev/sbom-digest: 7C1EB2358B426D9CC031B954A7CC402770A5DA0DB6CB83B507967055AB826A38`

```text
Name:             user-service-fddcb8687-dzjrh
Namespace:        stock-trading
Node:             desktop-worker/172.20.0.5
Status:           Running
Ready:            True
Image:            ttl.sh/stock-trading-d4ec05c5f73c@sha256:42e4b1f9ff6e20ec9ec3885614552472b680954f818456caa7aad6b675d36724
Events:
  Normal  Scheduled  Successfully assigned stock-trading/user-service-fddcb8687-dzjrh to desktop-worker
  Normal  Pulling    Pulling image "ttl.sh/stock-trading-d4ec05c5f73c@sha256:42e4b1f9ff6e20ec9ec3885614552472b680954f818456caa7aad6b675d36724"
  Normal  Pulled     Successfully pulled image in 2.188s
  Normal  Created    Created container: user-service
  Normal  Started    Started container user-service
```

## Runtime Logs
```text
2026/04/06 05:43:36 INFO SERVER START CONFIG config="{\"env\":\"local\",\"grpc\":{\"host\":\"0.0.0.0\",\"port\":19090},\"http\":{\"host\":\"0.0.0.0\",\"port\":18080},\"db\":{\"host\":\"127.0.0.1\",\"port\":3306,\"user\":\"root\",\"password\":\"\",\"name\":\"stock\"},\"auth\":{\"access_token_secret\":\"***\",\"access_token_ttl_minutes\":15,\"refresh_token_secret\":\"***\",\"refresh_token_ttl_minutes\":4320,\"issuer\":\"stock-trading-be\",\"audience\":\"stock-trading-clients\"},\"notification\":{\"kafka\":{\"brokers\":[\"localhost:29092\"],\"topic\":\"\",\"group_id\":\"email-service\"},\"email\":{\"provider\":\"noop\",\"smtp\":{\"host\":\"localhost\",\"port\":1025,\"username\":\"\",\"password\":\"\",\"from\":\"no-reply@example.com\",\"use_tls\":false},\"verification_url_base\":\"http://127.0.0.1:18080/users/verify?token=\"}},\"verification\":{\"token_ttl_hours\":24,\"resend_cooldown_seconds\":60}}"
2026/04/06 05:43:36 ERROR MYSQL UNRESPONSIVE error="dial tcp 127.0.0.1:3306: connect: connection refused"
2026/04/06 05:43:36 INFO NOTIFICATION SERVICE DISABLED reason="missing configuration"
2026/04/06 05:43:36 INFO SERVER STARTED
2026/04/06 05:43:36 INFO HTTP GATEWAY RUNNING addr=0.0.0.0:18080
2026/04/06 05:43:36 INFO GRPC SERVER RUNNING addr=[::]:19090
```

## Deployment Events
```text
LAST SEEN   TYPE     REASON              OBJECT                              MESSAGE
4m18s       Normal   ScalingReplicaSet   deployment/user-service             Scaled up replica set user-service-fddcb8687 from 0 to 1
4m16s       Normal   Scheduled           pod/user-service-fddcb8687-dzjrh    Successfully assigned stock-trading/user-service-fddcb8687-dzjrh to desktop-worker
4m16s       Normal   Pulling             pod/user-service-fddcb8687-dzjrh    Pulling image "ttl.sh/stock-trading-d4ec05c5f73c@sha256:42e4b1f9ff6e20ec9ec3885614552472b680954f818456caa7aad6b675d36724"
4m16s       Normal   SuccessfulCreate    replicaset/user-service-fddcb8687   Created pod: user-service-fddcb8687-dzjrh
4m13s       Normal   Pulled              pod/user-service-fddcb8687-dzjrh    Successfully pulled image "ttl.sh/stock-trading-d4ec05c5f73c@sha256:42e4b1f9ff6e20ec9ec3885614552472b680954f818456caa7aad6b675d36724" in 2.188s (2.188s including waiting). Image size: 16639135 bytes.
4m13s       Normal   Created             pod/user-service-fddcb8687-dzjrh    Created container: user-service
4m13s       Normal   Started             pod/user-service-fddcb8687-dzjrh    Started container user-service
```

## Kyverno Local Demo Policy
The local passing demo used an additional Kyverno policy dedicated to the temporary signed image on `ttl.sh`.

```text
ClusterPolicy: verify-local-demo-image
Rule: verify-local-demo-signature
Image reference match: ttl.sh/stock-trading-d4ec05c5f73c*
Status: Ready
```

## Notes
- This passing scenario was validated on the `docker-desktop` cluster, not the legacy `kind-devsecops` context.
- The image is hosted on `ttl.sh`, so it is intentionally temporary.
- The application starts successfully in local demo mode with Kubernetes admission verification enabled, while database access falls back because no in-cluster MySQL service is provisioned in this demo path.
