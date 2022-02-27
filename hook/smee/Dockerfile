FROM node:16.14.0 as build

WORKDIR /usr/src/app
COPY package*.json ./
COPY main.js ./

RUN npm install

FROM gcr.io/distroless/nodejs
COPY --from=build /usr/src/app /
EXPOSE 8080
CMD ["main.js"]