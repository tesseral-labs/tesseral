FROM nginx:1.29.0

COPY nginx.conf /etc/nginx/nginx.conf
COPY public /www/public

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]
