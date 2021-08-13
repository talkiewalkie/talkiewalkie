#!/bin/bash

clean() {
	rm -f "kanikojob*.yaml"
}

trap clean EXIT ERR

set -ex

SHA="$(git rev-parse HEAD)"
BUILD_NAME="back-build-context-$SHA.tar.gz"
BUILD_ARCHIVE="/tmp/$BUILD_NAME"
BLOB_NAME="gs://talkiewalkie-dev/kaniko/$BUILD_NAME"
JOB_FILE="kanikojob-$SHA.yaml"

pushd ..
tar -czf "$BUILD_ARCHIVE" \
  --exclude server/sqlboiler.toml \
  protos \
  server
popd

gsutil cp "$BUILD_ARCHIVE" gs://talkiewalkie-dev/kaniko/

cat <<EOF > $JOB_FILE
apiVersion: batch/v1
kind: Job
metadata:
  name: back-build-$SHA
spec:
  ttlSecondsAfterFinished: 120
  backoffLimit: 0
  template:
    spec:
      restartPolicy: Never
      serviceAccountName: kaniko-sa
      containers:
      - name: kaniko
        image: gcr.io/kaniko-project/executor:latest
        args: ["--dockerfile=server/Dockerfile",
               "--context=$BLOB_NAME",
               "--cache=true",
               "--destination=gcr.io/talkiewalkie-305117/talkiewalkie-back:$SHA"]
EOF

kubectl apply -f $JOB_FILE -n prod
