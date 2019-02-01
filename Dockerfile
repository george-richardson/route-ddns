FROM alpine:3.9
RUN apk add --no-cache libc6-compat ca-certificates
ADD route-ddns /route-ddns
CMD /route-ddns