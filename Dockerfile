FROM mhart/alpine-node:6.10.2

RUN apk add --no-cache bash git zip curl make

# install AWS_CLI (via python)
RUN apk -Uuv add groff less python py-pip && \
pip install awscli && \
apk --purge -v del py-pip && \
rm /var/cache/apk/*

WORKDIR /go/src/adv-caja-x-api/bin
#add serverless dependency
RUN npm install -g serverless@1.26.0

COPY ./package.json .
COPY ./package-lock.json .
RUN npm install
