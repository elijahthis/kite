# Build Stage
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production Stage (Nginx)
FROM nginx:alpine
# Copy the built React app
COPY --from=build /app/dist /usr/share/nginx/html
# Copy the custom Nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]