#!/usr/bin/env bash

cleanup() {
  rm -f Dockerfile
}
trap cleanup ERR EXIT

cat <<EOF >>Dockerfile
FROM nginx:alpine
RUN apk --no-cache add ca-certificates
WORKDIR /usr/share/nginx/html
COPY web/build/ .
COPY nginx /etc/nginx/
CMD ["nginx", "-g", "daemon off;"]
EOF

set -eEx

pushd web
yarn && yarn build
popd

sha="$(git rev-parse HEAD)"
docker build --platform linux/amd64 -t gcr.io/talkiewalkie-305117/talkiewalkie-front:${sha} -t talkiewalkie-front .
