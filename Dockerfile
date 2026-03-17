FROM alpine:latest

COPY generate-logs.sh /generate-logs.sh
RUN chmod +x /generate-logs.sh

# Default: 5 KB/s. Override with: docker run log-generator 10
ENTRYPOINT ["/generate-logs.sh"]
CMD ["5"]
