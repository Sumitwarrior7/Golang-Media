# Golang Media

Golang Media is a social media platform where users can create posts, view others' posts, add comments, and follow other users.

## Tech Stack

- **Frontend:** React, TailwindCSS
- **Backend:** Golang, SQL, Docker, PostgreSQL, Redis, Mailtrap

## Features

- **Containerized 3-tier architecture** for frontend, backend, and database services.
- **Fully customizable architecture** for third-party services by configuring database, caching, and mailing services.
- **JWT token-based authentication** with email verification and role-based authorization.
- **Paginated user feed and user management sections** for seamless browsing.
- **Search functionality** to find posts and other users efficiently.
- **Structured logging** for effective error handling and debugging.
- **Database seeding** to prepopulate with dummy data.
- **Makefile for automation** (running migrations, seeding DB, testing the application).
- **API rate limiting** to enhance security.
- **Unit tests and mock tests** to ensure robustness.

## Installation & Setup

### 1. Start Frontend Service
```sh
cd frontend
npm install
npm start
```

### 2. Start Backend Service
```sh
cd backend
go mod tidy
air
```

### 3. Start Third-Party Services (Database, Caching, Mailing)
Ensure Docker is installed and run:
```sh
cd backend
docker-compose up
```

## API Reference

For detailed API documentation, refer to the server API reference:

[Server API Reference](https://documenter.getpostman.com/view/41352298/2sAYQiBnYP)

## Contributing

We welcome contributions! To contribute:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any queries or feedback, feel free to reach out:
- **Email:** sumitheavydriver2017@gmail.com
- **GitHub Issues:** [Report an issue](https://github.com/golangmedia/issues)
- **Community Forum:** [Join the discussion](https://community.golangmedia.com)

---
Thank you for using **Golang Media**! 🚀