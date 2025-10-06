# # Stage 1: Build the frontend
# FROM node:20-alpine AS builder

# WORKDIR /app

# # Copy package.json and lock file first for caching
# COPY package*.json ./

# # Install dependencies
# RUN npm install

# # Copy source code
# COPY . .

# # Build the Vite app
# RUN npm run build

# # Stage 2: Serve the frontend using nginx
# FROM nginx:alpine

# # Copy build output to nginx html directory
# COPY --from=builder /app/dist /usr/share/nginx/html

# # Copy custom nginx config (optional)
# # COPY nginx.conf /etc/nginx/conf.d/default.conf

# # Expose port 5173 (Vite default) or 80 (nginx)
# EXPOSE 80

# CMD ["nginx", "-g", "daemon off;"]


# Use lightweight Node image
FROM node:20-alpine

WORKDIR /app

# Copy package.json and lock file
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy all source code
COPY . .

# Expose Vite dev server port
EXPOSE 5173

# Start Vite dev server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
