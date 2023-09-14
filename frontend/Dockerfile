# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:19.5.0-alpine3.16 AS build-stage
WORKDIR /app
# cache packages in layer
COPY package.json /app/package.json
COPY package-lock.json /app/package-lock.json
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci
# install
COPY . /app
RUN npm run build
ARG TARGETOS TARGETARCH
RUN echo $TARGETOS $TARGETARCH

# FROM alpine
# production stage
FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY resources/nginx.conf /etc/nginx/conf.d
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]