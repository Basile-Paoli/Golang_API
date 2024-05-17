FROM node:current-alpine3.19
WORKDIR /app
COPY package.json package.json
COPY package-lock.json package-lock.json
RUN npm ci
COPY . .
RUN npm run build
CMD ["npm","run","preview"]