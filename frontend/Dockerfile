FROM node:16-alpine

WORKDIR /app
COPY . ./

RUN npm config set registry https://registry.npmmirror.com && npm install

EXPOSE 3000

CMD ["npm", "start"]