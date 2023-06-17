<div align="center"><i>simple and open.</i></div>

#### 

This folder contains the source code for backend of the Primrose platform, the backend is designed to be *open*,
therefore, you can use it as your backend for your blog site without using the Primrose client.

Primrose's backend uses a customized authentication style called `Mengmeng` wherein tokens are 
signed two-ways using a database-stored user-token and the server-token, but in general, authentication is 
handled simply the same way as JWTs are handled through `Authorization` header.

MongoDB is used as  the database for this, it is expected that you should either:
1. Self-host MongoDB (recommended to use [`Docker`](https://hub.docker.com/_/mongo)).
2. Use MongoDB Atlas (recommended to use especially when you are unsure over how to self-host).

#### Installation

1. Configure the configurations needed by copying the `.env.example` to `.env` and completing the required configurations.
2. Build the Docker image: `docker build -t primrose-core .`
3. Run the Docker container: `docker run --name primrose-core -p 9053:9053 --env-file .env primrose-core`