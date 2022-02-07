FROM node:17.4-alpine3.14

WORKDIR /app

COPY package.json ./

RUN npm install

COPY . .

CMD [ "npm", "start" ]