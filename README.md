# singing-contest

- [x] Run something that serves
- [x] Use docker-compose
- [x] Add a database (postgresql) to a server
- [x] Implement a service: lists all contestants
- [ ] Add user authentication
- [ ] Add a service: receive voting result and show
- [ ] Add another service for votes: producer of events
- [ ] Build backend more as needed
- [ ] Build frontend
- [ ] Deploy

## are microservices the right solution for you?

No, because we don't have enough users to require any form of scaling yet. But, we are going to do it anyways -- for various reasons.

## Developer Notes

```sh
go run cmd/competition/main.go

# in another terminal
curl localhost:8000/
```

```sh
go mod vendor
docker-compose up --build

# in another terminal
curl localhost:8080/contestants
```
