# kube

This directory contains various utils for our GKE cluster. Services files are included in their respective codebases,
i.e. the deployment & service specs for the webapp are in `web/webapp.yaml`, for the server in `server/back.yaml`.

- `chisel.yaml` commands our [chisel](https://github.com/jpillora/chisel) tunneling for mobile development.
- `ingress.yaml` commands the Ingress rules (reverse-proxying).
- `certs.yaml` specs the domains we want managed certs for.

The missing file is `secrets.yaml` which has the production secrets for e.g. db credentials etc.

