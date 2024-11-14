FROM node:22.11.0-alpine
 
WORKDIR /usr/src/app
 
COPY package*.json ./
 
RUN npm install
 
COPY . .
 
EXPOSE 8080
 
CMD [ "npm", "start" ]