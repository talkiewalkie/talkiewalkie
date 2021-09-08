#!/bin/bash

clean() {
	rm -f kanikojob-*
}

trap clean EXIT ERR

set -ex

SHA="$(git rev-parse HEAD)"
BUILD_NAME="web-build-context-$SHA.tar.gz"
BUILD_ARCHIVE="/tmp/$BUILD_NAME"
BLOB_NAME="gs://talkiewalkie-dev/kaniko/$BUILD_NAME"
JOB_FILE="kanikojob-$SHA.yaml"

tar -czf "$BUILD_ARCHIVE" \
  --exclude node_modules \
  --exclude .next \
  --exclude .secrets \
  --exclude .env.local \
  --exclude .env.local.sample \
  .

gsutil cp "$BUILD_ARCHIVE" gs://talkiewalkie-dev/kaniko/

cat <<EOF > $JOB_FILE
apiVersion: batch/v1
kind: Job
metadata:
  name: web-build-$SHA
spec:
  ttlSecondsAfterFinished: 10
  backoffLimit: 0
  template:
    spec:
      restartPolicy: Never
      serviceAccountName: kaniko-sa
      containers:
      - name: kaniko
        image: gcr.io/kaniko-project/executor:latest
        # build-args might seem weird, but it works out:
        # https://github.com/GoogleContainerTools/kaniko/issues/713#issuecomment-589472822
        args: ["--dockerfile=Dockerfile",
               "--context=$BLOB_NAME",
               "--cache=true",
               "--destination=gcr.io/talkiewalkie-305117/talkiewalkie-front:$SHA",
               "--build-arg=FIREBASE_PRIVATE_KEY",
               "--build-arg=FIREBASE_CLIENT_EMAIL",
               "--build-arg=COOKIE_SECRET_CURRENT",
               "--build-arg=COOKIE_SECRET_PREVIOUS"]
        env:
          - name: FIREBASE_PRIVATE_KEY
            valueFrom:
              secretKeyRef:
                name: mysecrets
                key: firebase_api_key
          - name: FIREBASE_CLIENT_EMAIL
            valueFrom:
              secretKeyRef:
                name: mysecrets
                key: firebase_client_email
          - name: COOKIE_SECRET_CURRENT
            valueFrom:
              secretKeyRef:
                name: mysecrets
                key: cookie_secret_current
          - name: COOKIE_SECRET_PREVIOUS
            valueFrom:
              secretKeyRef:
                name: mysecrets
                key: cookie_secret_previous
EOF

kubectl apply -f $JOB_FILE -n prod
