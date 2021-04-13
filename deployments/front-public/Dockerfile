FROM nginx:alpine

# followed this guide
# https://geko.cloud/nginx-and-ssl-with-certbot-in-docker-alpine/
RUN apk add certbot certbot-nginx && mkdir -p /etc/letsencrypt
# docker build -t mynginx --no-cache .
# once build run
# docker run -v $(pwd)/letsencrypt:/etc/letsencrypt --name nginx -ti --network host --rm mynginx certbot --nginx -d talkiewalkie.app -d www.talkiewalkie.app -d dev.talkiewalkie.app -d api.talkiewalkie.app -d tunnel.talkiewalkie.app
# chisel setup: on laptop: chisel client -v https://tunnel.talkiewalkie.app R:8081:127.0.0.1:8080
#               on server: chisel server -v -p 8002 --reverse


COPY nginx /etc/nginx/

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
