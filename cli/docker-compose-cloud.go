package cli

const composeCloud = `
version: "3.7"

services:

  pubsub:
    image: rabbitmq:3-management
    ports:
      - 5672:5672
      - 15672:15672
    labels:
      NAME: rabbitmq
    volumes:
      - ${RABBITMQ_DATA_VOLUME}:/var/lib/rabbitmq
    networks: 
      - pubsub
    healthcheck:
      test: [ "CMD", "rabbitmqctl", "status" ]
      interval: 60s
      timeout: 30s
    deploy:
      placement:
          constraints:
            - node.role == manager

  db:
    image: postgres
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    ports:
      - 5432:5432
    volumes:
      - ${POSTGRES_DATA_VOLUME}:/var/lib/postgresql/data
    networks: 
      - db
    deploy:
      placement:
          constraints:
            - node.role == manager
  
  worker:
    init: true
    image: $KENZA_CONTAINER_REGISTRY/worker:$KENZA_VERSION
    command: |
      /kenza/worker -logfile_dir=$KENZA_JOB_LOGS_CONTAINER_PATH
    networks: 
      - pubsub
    depends_on:
      - api
      - pubsub
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock # use host's docker
      - ${KENZA_JOB_LOGS_HOST_PATH}:${KENZA_JOB_LOGS_CONTAINER_PATH}:rw
    secrets:
      - api_key
    deploy:
      placement:
          constraints:
            - node.role == manager

  progress:
    init: true
    image: $KENZA_CONTAINER_REGISTRY/progress:$KENZA_VERSION
    command: /kenza/progress
    deploy:
      placement:
          constraints:
            - node.role == manager
    networks: 
      - pubsub
    depends_on:
      - api
      - pubsub
    secrets:
      - api_key

  api:
    init: true
    image: $KENZA_CONTAINER_REGISTRY/api:$KENZA_VERSION
    command: |
      /kenza/apid
      -db_name=$POSTGRES_DB
      -db_host=$POSTGRES_HOST
      -db_user=$POSTGRES_USER
      -db_pass=$POSTGRES_PASSWORD
      -logfile_dir=$KENZA_JOB_LOGS_CONTAINER_PATH
      -api_host=$KENZA_API_HOST
      -api_port=$KENZA_API_PORT
    ports:
      - 8080:8080
    networks:
      - db
      - pubsub
    depends_on:
      - db
      - pubsub
    volumes:
      - ${KENZA_JOB_LOGS_HOST_PATH}:${KENZA_JOB_LOGS_CONTAINER_PATH}:ro
    secrets:
      - api_key
      - github_webhook_secret
      
    deploy:
      placement:
          constraints:
            - node.role == manager
  
  scheduler:
    init: true
    image: $KENZA_CONTAINER_REGISTRY/scheduler:$KENZA_VERSION
    command: /kenza/scheduler 
      
    networks:
      - pubsub
    depends_on:
      - api
      - pubsub
    secrets:
      - api_key
    deploy:
      placement:
          constraints:
            - node.role == manager

  web:
    init: true
    image: $KENZA_CONTAINER_REGISTRY/web:$KENZA_VERSION
    ports:
      - 80:80
    depends_on:
      - api
    deploy:
      placement:
        constraints:
          - node.role == manager

networks: 
  db:
    name: kenza_network_db
  pubsub:
    name: kenza_network_pubsub

secrets:
  api_key:
    file: ./api_key.secret
  github_webhook_secret:
    file: ./github_webhook_secret.secret
`
