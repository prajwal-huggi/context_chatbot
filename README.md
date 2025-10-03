Go Backend Setup
to initialize the go: go mod init github.com/prajwal-huggi/backend_go
from the link(go lang clean)-> https://github.com/ilyakaznacheev/cleanenv

go get -u github.com/ilyakaznacheev/cleanenv

go run cmd/server/main.go -config config/local.yaml
-config is the flag which is must

https://github.com/go-playground/validator
go get github.com/go-playground/validator/v10
The above command is used to validate the request which is sent by the user.

installing the go sqlite driver
https://github.com/mattn/go-sqlite3
go get github.com/mattn/go-sqlite3

github.com/joho/godotenv
The above github is used to load the .env file in the go project.(Usage is present in the backend/cmd/server/main.go)

FRONTEND
# React + Vite

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Babel](https://babeljs.io/) (or [oxc](https://oxc.rs) when used in [rolldown-vite](https://vite.dev/guide/rolldown)) for Fast Refresh
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/) for Fast Refresh

## React Compiler

The React Compiler is not enabled on this template because of its impact on dev & build performances. To add it, see [this documentation](https://react.dev/learn/react-compiler/installation).

## Expanding the ESLint configuration

If you are developing a production application, we recommend using TypeScript with type-aware lint rules enabled. Check out the [TS template](https://github.com/vitejs/vite/tree/main/packages/create-vite/template-react-ts) for information on how to integrate TypeScript and [`typescript-eslint`](https://typescript-eslint.io) in your project.
