FROM gcr.io/distroless/base
COPY api /api

EXPOSE 30018

ENTRYPOINT ["/api"]