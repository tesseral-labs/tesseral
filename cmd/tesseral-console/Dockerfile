FROM nginx:1.29.0

COPY nginx.conf /etc/nginx/nginx.conf
COPY public /www/public
COPY entrypoint /entrypoint

RUN chmod +x /entrypoint

EXPOSE 80

ENTRYPOINT ["/entrypoint"]
