# =========================
# Build stage
# =========================
FROM node:20-slim AS build

WORKDIR /app

# 先复制依赖文件（利用缓存）
COPY package*.json ./

# ❗必须安装全部依赖（不要 production）
RUN npm ci

# 再复制源码
COPY . .

# 构建
RUN npm run build


# =========================
# Production stage
# =========================
FROM nginx:alpine

# 复制构建产物
COPY --from=build /app/dist /usr/share/nginx/html

# nginx 配置
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]