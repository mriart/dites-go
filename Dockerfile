FROM ubuntu:latest


COPY . ./

EXPOSE 8080
CMD ["./dites-go"]
