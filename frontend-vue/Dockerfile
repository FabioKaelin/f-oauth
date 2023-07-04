# Multi-stage
# 1) Node image for building frontend assets
# 2) nginx stage to serve frontend assets

# Name the node stage "builder"
FROM node:lts-alpine AS builder
# Set working directory
WORKDIR /app
# Copy all files from current directory to working dir in image
COPY . .
# install node modules and build assets
RUN npm install
RUN npm run build

# # nginx state for serving content
# FROM nginx:alpine
# # Set working directory to nginx asset directory
# WORKDIR /usr/share/nginx/html
# # Remove default nginx static assets
# RUN rm -rf ./*
# # Copy static assets from builder stage
# COPY --from=builder /app/dist .
# # Containers run nginx with global directives and daemon off
# ENTRYPOINT ["nginx", "-g", "daemon off;"]

FROM nginx:1.19

RUN rm -f /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d

RUN chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/conf.d

RUN touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

USER nginx

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 8080