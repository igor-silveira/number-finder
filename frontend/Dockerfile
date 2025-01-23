FROM node:19-alpine AS builder

WORKDIR /app

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

COPY package*.json ./

RUN npm install

COPY .. .

RUN npm run build

FROM node:19-alpine

WORKDIR /app

COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

RUN npm install --production

ENV NODE_ENV=production

CMD ["npm", "start"]
