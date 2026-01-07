# Demo Evidence - Supply Chain Security

## Image Digest
Image: ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a

## Cosign Verify (signature)
```

Verification for ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a --
The following checks were performed on each of these signatures:
  - The cosign claims were validated
  - Existence of the claims in the transparency log was verified offline
  - The signatures were verified against the specified public key

[{"critical":{"identity":{"docker-reference":"ghcr.io/sinhnguyen1411/stock-trading/user-service"},"image":{"docker-manifest-digest":"sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a"},"type":"cosign container image signature"},"optional":{"Bundle":{"SignedEntryTimestamp":"MEQCIHbimnCpdNjfiKJrXfQJKJAsTgq2kSFL4L3Z4Q124owRAiAclExXoo1lHVdoxlG8qO4vqf+3rV6kHqw1BSlV3WidrA==","Payload":{"body":"eyJhcGlWZXJzaW9uIjoiMC4wLjEiLCJraW5kIjoiaGFzaGVkcmVrb3JkIiwic3BlYyI6eyJkYXRhIjp7Imhhc2giOnsiYWxnb3JpdGhtIjoic2hhMjU2IiwidmFsdWUiOiI5MWViN2ZkZWI5NThkMTgzMTJjMGYzZmE4OWU5OWEyNDMwNDllNjM3MTQ5MzE5NjlhYjRkZmI4YmVlZjA2Y2UyIn19LCJzaWduYXR1cmUiOnsiY29udGVudCI6Ik1FUUNJRURhRDJ5QzZwRXpkRUthM1ZHc3lZUGV2SVp6Z2RiQ3dKc0hMNW9Ud2k1aEFpQXo3aTVrbE5BeU4rNnM3QmVab0F3NzQ2QXR0cDRoTTFEYlhuenVhekMvNHc9PSIsInB1YmxpY0tleSI6eyJjb250ZW50IjoiTFMwdExTMUNSVWRKVGlCUVZVSk1TVU1nUzBWWkxTMHRMUzBLVFVacmQwVjNXVWhMYjFwSmVtb3dRMEZSV1VsTGIxcEplbW93UkVGUlkwUlJaMEZGVG1ZeVVrbFdUSGhSYWxoNmNuY3ZOVGhDVWk5ME9YbEdSbGQxTkFvMlEybEdSME5XYUVFNVJtdDVaa05HUnpKS2JFdFRRbmhMTjJGSGExbzNiMmRPWTFsM1JIZEJZbk5YWTNWak0xUkVSVEl4SzBWd1MyZG5QVDBLTFMwdExTMUZUa1FnVUZWQ1RFbERJRXRGV1MwdExTMHRDZz09In19fX0=","integratedTime":1767758551,"logIndex":799068178,"logID":"c0d23d6ad406973f9559f3ba2d1ca01f84147d8ffc5b8445c224f98b9591801d"}}}}]
```

## Cosign Verify (attestation/provenance)
```

Verification for ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a --
The following checks were performed on each of these signatures:
  - The cosign claims were validated
  - Existence of the claims in the transparency log was verified offline
  - The signatures were verified against the specified public key
{"payloadType":"application/vnd.in-toto+json","payload":"eyJfdHlwZSI6Imh0dHBzOi8vaW4tdG90by5pby9TdGF0ZW1lbnQvdjAuMSIsInByZWRpY2F0ZVR5cGUiOiJodHRwczovL3Nsc2EuZGV2L3Byb3ZlbmFuY2UvdjAuMiIsInN1YmplY3QiOlt7Im5hbWUiOiJnaGNyLmlvL3NpbmhuZ3V5ZW4xNDExL3N0b2NrLXRyYWRpbmcvdXNlci1zZXJ2aWNlIiwiZGlnZXN0Ijp7InNoYTI1NiI6ImFmMDNkZjIyYTlkYzUwNTgxMDlmYTA4OWRiYTI0ZWFjYTE0NWY2ZDM2M2FhMTEzYmIzMGMyYmZkMjQzYTc3MGEifX1dLCJwcmVkaWNhdGUiOnsiYnVpbGRUeXBlIjoiaHR0cHM6Ly9zbHNhLmRldi9wcm92ZW5hbmNlL3YwLjIiLCJidWlsZGVyIjp7ImlkIjoibG9jYWwtZGVtbyJ9LCJpbnZvY2F0aW9uIjp7ImNvbmZpZ1NvdXJjZSI6eyJkaWdlc3QiOnsiZ2l0Q29tbWl0IjoiZGVtbyJ9LCJ1cmkiOiJodHRwczovL2dpdGh1Yi5jb20vc2luaG5ndXllbjE0MTEvc3RvY2stdHJhZGluZy1iZSJ9fX19","signatures":[{"keyid":"","sig":"MEUCID10wvA5ddYrSyUyb0PBzeDrhqbt5akK4UR3j2InzxwZAiEAng60nrLxEK2qEB8Fl3gEbuW6/fcmK3tM2xnYFNw8tmE="}]}
```

## Kyverno Admission Deny Evidence (CVE gate)
```
LAST SEEN   TYPE     REASON              OBJECT                               MESSAGE
18m         Normal   SandboxChanged      pod/user-service-c9b5df5d9-kk67h     Pod sandbox changed, it will be killed and re-created.
18m         Normal   Started             pod/user-service-c9b5df5d9-kk67h     Started container user-service
18m         Normal   Created             pod/user-service-c9b5df5d9-kk67h     Created container user-service
18m         Normal   Pulled              pod/user-service-c9b5df5d9-kk67h     Container image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:e67b03ed25ed0b7a6c7aa3c8530d997a3396913581128a91083e7b31270bd121" already present on machine
8m42s       Normal   SuccessfulCreate    replicaset/user-service-6cdb4c544c   Created pod: user-service-6cdb4c544c-4d942
8m42s       Normal   ScalingReplicaSet   deployment/user-service              Scaled up replica set user-service-6cdb4c544c to 1
8m41s       Normal   Scheduled           pod/user-service-6cdb4c544c-4d942    Successfully assigned stock-trading/user-service-6cdb4c544c-4d942 to devsecops-control-plane
8m41s       Normal   SuccessfulCreate    replicaset/user-service-6cdb4c544c   Created pod: user-service-6cdb4c544c-4jzbf
8m40s       Normal   Scheduled           pod/user-service-6cdb4c544c-4jzbf    Successfully assigned stock-trading/user-service-6cdb4c544c-4jzbf to devsecops-control-plane
8m41s       Normal   SuccessfulCreate    replicaset/user-service-c9b5df5d9    Created pod: user-service-c9b5df5d9-8jnr5
8m41s       Normal   Killing             pod/user-service-c9b5df5d9-kk67h     Stopping container user-service
8m41s       Normal   Pulling             pod/user-service-6cdb4c544c-4d942    Pulling image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a"
8m40s       Normal   Scheduled           pod/user-service-c9b5df5d9-8jnr5     Successfully assigned stock-trading/user-service-c9b5df5d9-8jnr5 to devsecops-control-plane
8m40s       Normal   Pulling             pod/user-service-6cdb4c544c-4jzbf    Pulling image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a"
8m40s       Normal   Started             pod/user-service-c9b5df5d9-8jnr5     Started container user-service
8m40s       Normal   Created             pod/user-service-c9b5df5d9-8jnr5     Created container user-service
8m40s       Normal   Pulled              pod/user-service-c9b5df5d9-8jnr5     Container image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:e67b03ed25ed0b7a6c7aa3c8530d997a3396913581128a91083e7b31270bd121" already present on machine
8m36s       Normal   Pulled              pod/user-service-6cdb4c544c-4d942    Successfully pulled image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a" in 4.79s (4.79s including waiting). Image size: 16525352 bytes.
8m36s       Normal   Started             pod/user-service-6cdb4c544c-4d942    Started container user-service
8m36s       Normal   Created             pod/user-service-6cdb4c544c-4d942    Created container user-service
8m36s       Normal   Killing             pod/user-service-6cdb4c544c-4d942    Stopping container user-service
8m35s       Normal   Started             pod/user-service-6cdb4c544c-4jzbf    Started container user-service
8m35s       Normal   Pulled              pod/user-service-6cdb4c544c-4jzbf    Successfully pulled image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a" in 884ms (4.947s including waiting). Image size: 16525352 bytes.
8m35s       Normal   Created             pod/user-service-6cdb4c544c-4jzbf    Created container user-service
8m20s       Normal   ScalingReplicaSet   deployment/user-service              Scaled down replica set user-service-c9b5df5d9 to 0 from 1
8m20s       Normal   SuccessfulDelete    replicaset/user-service-c9b5df5d9    Deleted pod: user-service-c9b5df5d9-8jnr5
8m20s       Normal   Killing             pod/user-service-c9b5df5d9-8jnr5     Stopping container user-service
7m15s       Normal   Pulled              pod/user-service-6cdb4c544c-ctp2n    Container image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a" already present on machine
7m15s       Normal   Killing             pod/user-service-6cdb4c544c-4jzbf    Stopping container user-service
7m14s       Normal   Scheduled           pod/user-service-6cdb4c544c-ctp2n    Successfully assigned stock-trading/user-service-6cdb4c544c-ctp2n to devsecops-control-plane
7m15s       Normal   SuccessfulCreate    replicaset/user-service-6cdb4c544c   Created pod: user-service-6cdb4c544c-ctp2n
7m15s       Normal   Created             pod/user-service-6cdb4c544c-ctp2n    Created container user-service
7m15s       Normal   Started             pod/user-service-6cdb4c544c-ctp2n    Started container user-service
```

## Current Pod Status
```
NAME                            READY   STATUS    RESTARTS   AGE
user-service-6cdb4c544c-ctp2n   1/1     Running   0          7m16s
```

## Kyverno Admission Deny Evidence (latest)
```
36s         Warning   FailedCreate        replicaset/user-service-87448b5f6    Error creating: admission webhook "validate.kyverno.svc-fail" denied the request: ...
36s         Warning   FailedCreate        replicaset/user-service-87448b5f6    Error creating: admission webhook "validate.kyverno.svc-fail" denied the request: ...
36s         Warning   FailedCreate        replicaset/user-service-87448b5f6    Error creating: admission webhook "validate.kyverno.svc-fail" denied the request: ...
36s         Warning   FailedCreate        replicaset/user-service-87448b5f6    Error creating: admission webhook "validate.kyverno.svc-fail" denied the request: ...
36s         Warning   FailedCreate        replicaset/user-service-87448b5f6    Error creating: admission webhook "validate.kyverno.svc-fail" denied the request: ...
36s         Normal    ScalingReplicaSet   deployment/user-service              Scaled up replica set user-service-87448b5f6 to 1
36s         Normal    Created             pod/user-service-6cdb4c544c-jw287    Created container user-service
36s         Normal    Pulled              pod/user-service-6cdb4c544c-jw287    Container image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a" already present on machine
35s         Normal    Scheduled           pod/user-service-6cdb4c544c-jw287    Successfully assigned stock-trading/user-service-6cdb4c544c-jw287 to devsecops-control-plane
36s         Normal    SuccessfulCreate    replicaset/user-service-6cdb4c544c   Created pod: user-service-6cdb4c544c-jw287
35s         Normal    Started             pod/user-service-6cdb4c544c-jw287    Started container user-service
35s         Warning   FailedCreate        replicaset/user-service-87448b5f6    Error creating: admission webhook "validate.kyverno.svc-fail" denied the request: ...
26s         Warning   FailedCreate        replicaset/user-service-87448b5f6    (combined from similar events): Error creating: admission webhook "validate.kyverno.svc-fail" denied the request: ...
33s         Normal    Pulled              pod/user-service-6cdb4c544c-dlngm    Container image "ghcr.io/sinhnguyen1411/stock-trading/user-service@sha256:af03df22a9dc5058109fa089dba24eaca145f6d363aa113bb30c2bfd243a770a" already present on machine
33s         Normal    Scheduled           pod/user-service-6cdb4c544c-dlngm    Successfully assigned stock-trading/user-service-6cdb4c544c-dlngm to devsecops-control-plane
33s         Normal    Created             pod/user-service-6cdb4c544c-dlngm    Created container user-service
33s         Normal    Started             pod/user-service-6cdb4c544c-dlngm    Started container user-service
33s         Normal    Killing             pod/user-service-6cdb4c544c-jw287    Stopping container user-service
33s         Normal    SuccessfulCreate    replicaset/user-service-6cdb4c544c   Created pod: user-service-6cdb4c544c-dlngm
23s         Normal    ScalingReplicaSet   deployment/user-service              Scaled down replica set user-service-87448b5f6 to 0 from 1
```
