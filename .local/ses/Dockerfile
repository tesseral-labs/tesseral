FROM node:22.11.0

WORKDIR /root

RUN mkdir -p /output

RUN npm install aws-ses-v2-local -g

EXPOSE 8005

ENTRYPOINT ["aws-ses-v2-local", "--port", "8005", "--host", "0.0.0.0"]