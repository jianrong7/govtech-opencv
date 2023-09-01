# govtech-opencv

## Demo (<https://opencv.jianrong.dev/>)

[Demo](https://opencv.jianrong.dev/)

[Postman API](https://www.postman.com/solar-sunset-754254/workspace/govtech-opencv/collection/16590827-d0a7429c-a582-45b7-9d75-c4de09b6d820?action=share&creator=16590827&active-environment=16590827-6621e113-a1ac-4a56-b70c-2e6cb57d21a9)

Hosted on render.com's free tier.

## Local installation

1. Clone the repository.
2. Start up Docker.
3. Run `docker-compose up -d` to start up the containers.
4. Your app is now running on `http://localhost:3000`.

## Testing

1. Go to `/router`.
2. Run `go test`.

`api_test.go` contains the unit tests.

## Tech stack

- Language: Golang

As per the requirement.

- Framework: Fiber

I have previously worked with Gin, one of the most popular frameworks for golang. This time, I wanted to try something new. I watched several videos about Fiber and thought it would be a good time to try this framework out. I liked how it is very similar to express so it offered a lot of flexibility while also came with many useful middlewares out of the box.

- Database: PostgreSQL

As per the requirement.

- ORM: GORM

I have used this ORM before and I quite liked it. It is popular and I could write raw SQL with it as well if I was unhappy with the generated queries.

- Containerization: Docker

Makes it easier for deployment.
