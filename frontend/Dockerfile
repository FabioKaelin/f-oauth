# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:lts AS build
WORKDIR /app
# cache packages in layer
COPY package.json /app/package.json
COPY package-lock.json /app/package-lock.json
RUN echo -n "NPM Version:" ; npm --version
RUN echo -n "Node Version:" ; node --version
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci
# install
COPY . /app
RUN npm run build
ARG TARGETOS TARGETARCH
RUN echo $TARGETOS $TARGETARCH

# production stage
FROM nginx:stable-alpine as production
LABEL org.opencontainers.image.source="https://github.com/FabioKaelin/f-oauth"
LABEL org.opencontainers.image.authors="FabioKaelin"
LABEL org.opencontainers.image.title="F-OAuth-Frontend"
COPY --from=build /app/dist /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY resources/nginx.conf /etc/nginx/conf.d
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
