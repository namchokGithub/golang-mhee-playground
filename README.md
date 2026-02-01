# Proundmhee

![Proundmhee](https://img.shields.io/badge/VocabBunny-🧸%20playground-DC8AFF)
![Version](https://img.shields.io/badge/version-1.0.0-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)

สนามเด็กเล่นโกลังหมี่ [Proundmhee]() ที่พัฒนาโดย [Namchok Singhachai]()

## Structure

<details>
  <summary>คลิกเพื่อดูโครงสร้างโปรเจ็ค</summary>

```text
proundmhee/
├── go.mod
├── cmd/
│   └── main/
│       └── main.go
└── internal/
    ├── app/
    │   ├── server.go
    │   └── routes.go
    ├── infra/
    │   ├── di/
    │   │   ├── contracts.go
    │   │   └── deps.go
    │   ├── logger/
    │   │   ├── gin.go
    │   │   └── logger.go
    ├── shared/
    │   ├── response.go
    │   └── middleware.go
    └── modules/
        ├── vat/
        │   ├── testing/
        │   │   └── handler_test.go
        │   │   └── service_test.go
        │   ├── handler.go
        │   ├── routes.go
        │   └── service.go
        ├── rsa/
        │   └── ...
        ├── generate_code/
        │   └── ...
        |── refundable_date/
        │   └── ...
        |── schemaground/
        │   └── ...
        │
        └── ... (Another modules)
```

</details>

## Project Setup

#### 1. สร้างโปรเจ็ค

`mkdir {project_name} && cd {project_name}`

> `{project_name}` เลือกชื่อโปรเจ็ค

#### 2. สร้างไฟล์ go.mod

`go mod init {project_name}`

#### 3. ติดตั้ง gin

`go get github.com/gin-gonic/gin`

#### 4. Implement modules

> โดยทำตามโครงสร้างของ [Gin](https://gin-gonic.com/en/docs/) และสร้าง modules ต่างๆ

#### 5. Run & test

`go run ./cmd/{project_name}`

`curl -s localhost:8080/health`

## Others

#### Testing

- รันทั้งหมด
  `go test ./...`

- เปิด coverage
  `go test ./... -cover`

- ล่า race condition
  `go test ./... -race`

</br>

---

_© 2026 Proundmhee. Released under the MIT License._
