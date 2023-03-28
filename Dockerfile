FROM betashepherd/golang:1.20-alpine3.17 as build
WORKDIR /build
ADD . .
RUN make build

FROM betashepherd/alpine:3.17
WORKDIR /opt/apps
COPY --from=build /build/chatgpt-web chatgpt-web
COPY --from=build /build/static static
COPY --from=build /build/resources resources