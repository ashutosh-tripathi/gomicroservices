FROM alpine:latest
RUN mkdir /app
COPY mail-service /app
CMD ["/app/mail-service"]