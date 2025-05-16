FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app .

FROM alpine:3.21.3 as production
LABEL org.opencontainers.image.source="https://github.com/FabioKaelin/f-oauth"
LABEL org.opencontainers.image.authors="FabioKaelin"
LABEL org.opencontainers.image.title="F-OAuth-Backend"
COPY --from=build /out/app /bin/app
CMD ["/bin/app"]
