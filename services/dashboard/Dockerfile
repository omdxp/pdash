FROM node:16.9.0 AS build
WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

ARG VITE_CUSTOMERS_URL
ARG VITE_ORDERS_URL
ARG VITE_SUPPLIERS_URL
ARG VITE_AUTH_URL

ENV VITE_CUSTOMERS_URL $VITE_CUSTOMERS_URL
ENV VITE_ORDERS_URL $VITE_ORDERS_URL
ENV VITE_SUPPLIERS_URL $VITE_SUPPLIERS_URL
ENV VITE_AUTH_URL $VITE_AUTH_URL

# install pnpm 
RUN npm install -g pnpm

# install dependencies
COPY package.json ./
COPY pnpm-lock.yaml ./
RUN pnpm install

# add app
COPY . ./
RUN pnpm run build

# host in nginx
FROM nginx:1.23.1-alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build /app/dist /dashboard/build