FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY multiPoolIncome /app/
ENTRYPOINT ["/app/multiPoolIncome"]
EXPOSE 8000