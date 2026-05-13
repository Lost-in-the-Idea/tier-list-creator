FROM node:lts-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
EXPOSE 4200
CMD ["npx", "ng", "serve", "--host", "0.0.0.0", "--poll", "500"]
