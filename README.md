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
