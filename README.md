# singing-contest

- [x] Run something that serves
- [ ] Use docker-compose
- [ ] Add a database (postgresql) to a server
- [ ] Implement a service that
  - lists all contestants
- [ ] Add user authentication
- [ ] Add another service for votes: producer of events
- [ ] Build backend
- [ ] Build frontend
- [ ] Deploy

## are microservices the right solution for you?

No, because we don't have enough users to require any form of scaling yet. But, we are going to do it anyways -- for various reasons.

## Developer Notes

```
go run cmd/competition/main.go
```
