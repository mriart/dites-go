FROM ubuntu:latest


COPY . ./
RUN mkdir -p /etc/ssl/certs/
COPY ./ca-certificates.crt /etc/ssl/certs/
ENV GEMINI_API_KEY=${GEMINI_API_KEY}

EXPOSE 8080
CMD ["./dites-go"]
