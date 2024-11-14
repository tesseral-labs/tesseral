FROM node:22.11.0-alpine
 
WORKDIR /usr/src/app
 
COPY package*.json ./
 
RUN npm install
 
COPY . .
 
EXPOSE 8082
 
CMD [ "npm", "start" ]