#!/usr/bin/env bash

cleanup () {
  rm -f Dockerfile
}

trap cleanup ERR EXIT

cat << EOF >> Dockerfile
FROM nginx:alpine
RUN apk --no-cache add ca-certificates
WORKDIR /usr/share/nginx/html
COPY build/ .
#COPY nginx /etc/nginx/
CMD ["nginx", "-g", "daemon off;"]
EOF

set -eEx

yarn && yarn build

docker build --platform linux/amd64 -t gcr.io/talkiewalkie-305117/talkiewalkie-front:1 .
