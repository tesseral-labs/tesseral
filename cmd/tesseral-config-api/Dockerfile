FROM nginx:1.29.0

COPY nginx.conf.template /etc/nginx/nginx.conf.template
COPY entrypoint /entrypoint

RUN chmod +x /entrypoint

EXPOSE 80

ENTRYPOINT ["/entrypoint"]
