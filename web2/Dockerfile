FROM node:14.15.4-alpine

WORKDIR /usr/app
# TODO: currently including .env.local in the build process :/
COPY . .
RUN yarn
RUN --mount=type=secret,id=mysecret,dst=.env.local yarn build

CMD [ "yarn", "start" ]