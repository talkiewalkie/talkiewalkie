FROM nginx:alpine

# followed this guide
# https://geko.cloud/nginx-and-ssl-with-certbot-in-docker-alpine/
RUN apk add certbot certbot-nginx && mkdir -p /etc/letsencrypt

COPY nginx /etc/nginx/

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
