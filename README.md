# gin-auth with mongoDB, support both REST API and GraphQL

This project uses the Gin framework to implement authentication and supports both REST API and GraphQL.

## Compatibility with REST API and GraphQL

This framework is designed to seamlessly support both RESTful APIs and GraphQL, allowing developers to choose the best approach for their application needs. Here’s how it achieves this compatibility:

1. **Unified Controller Logic**: 
   - The project structure includes a dedicated `controllers` directory where the logic for handling requests is implemented. This allows for a clear separation of concerns, making it easy to manage both REST and GraphQL endpoints within the same application.

2. **Routing**:
   - The `routes` directory defines the routing for RESTful endpoints, while the GraphQL endpoints are handled through a dedicated GraphQL server setup. This allows both types of requests to coexist without conflict.

3. **GraphQL Resolvers**:
   - The `graph` directory contains resolvers that handle GraphQL queries and mutations. These resolvers can interact with the same underlying services and models used by the REST API, ensuring that business logic is reused effectively.

4. **Shared Models**:
   - The `models` directory defines data structures that are used by both the REST API and GraphQL. This promotes consistency in data handling and reduces duplication of code.

5. **Middleware Support**:
   - Middleware functions defined in the `middlewares` directory can be applied to both REST and GraphQL routes. This allows for shared functionality such as authentication, logging, and error handling across both types of requests.

6. **Flexible Response Handling**:
   - The framework includes utilities for formatting responses, which can be adapted for both JSON responses (for REST) and GraphQL responses. This ensures that clients receive data in the expected format regardless of the endpoint type. Refer to `utils\response\response.go`

7. **Configuration**:
   - The application can be configured to enable or disable either REST or GraphQL features as needed, providing flexibility for different use cases.

By leveraging these design principles, the gin-auth framework provides a robust solution for building applications that require both RESTful and GraphQL APIs, allowing developers to choose the best approach for their specific requirements.



## Prerequisites

1. GOLANG Environment

   My Environment is `go version go1.23.2 windows/amd64`

   Download: https://go.dev/dl/

2. PostgreSQL

   Installation Guide: https://blog.csdn.net/weixin_68256171/article/details/132337173

3. Redis

4. minIO

   K8s: https://www.minio.org.cn/docs/minio/kubernetes/upstream/index.html
   https://blog.csdn.net/lichanggu/article/details/116270330
   Go Client API: https://min.io/docs/minio/linux/developers/go/API.html


## Run Project

1. Clone this repo

   ```
   git clone xxx
   cd gin-auth
   ```

2.  Install the dependencies

   ``` 
   go mod download
   ```

3. Setup environment variables

   create a `.env` file in project root directory, you can refer to `.env.example` or 

   ```shell
   # postgresql
   POSTGRES_USER=YOUR_POSTGRESQL_USERNAME
   POSTGRES_PASSWORD=YOUR_POSTGRESQL_PASSWORD
   POSTGRES_DB=YOUR_POSTGRESQL_DATABASE_NAME
   POSTGRES_HOST=YOUR_POSTGRESQL_HOST
   POSTGRES_PORT=YOUR_POSTGRESQL_PORT # default is 5432
   
   # redis
   REDIS_HOST=YOUR_REDIS_HOST
   REDIS_PORT=YOUR_REDIS_PORT # default is 6379
   REDIS_PASSWORD=YOUR_REDIS_PASSWORD
   REDIS_DB=YOUR_REDIS_DB # default is 0
   
   # minio
   MINIO_URL=YOUR_MINIO_URL # URL CANNOT contain http:// or https://
   MINIO_ACCESS_KEY_ID=YOUR_MINIO_ACCESS_KEY_ID
   MINIO_ACCESS_KEY_SECRET=YOUR_MINIO_ACCESS_KEY_SECRET
   MINIO_PUBLIC_BUCKET_NAME=YOUR_MINIO_PUBLIC_BUCKET_NAME # public
   MINIO_ROOT_USER=YOUR_MINIO_ROOT_USER
   MINIO_ROOT_PASSWORD=YOUR_MINIO_ROOT_PASSWORD
   
   # backend
   PREFIX=http://
   HOST=YOUR_HOST
   PORT=YOUR_PORT # default is 8080
   BACKEND_URL=${PREFIX}${HOST}:${PORT}
   
   # secret
   ARGON2_SALT=YOUR_ARGON2_SALT
   
   # smtp config
   SMTP_SERVER=YOUR_SMTP_SERVER # example: smtp.163.com
   SMTP_USER=YOUR_SMTP_USER
   SMTP_PASS=YOUR_SMTP_PASS
   SMTP_PORT=YOUR_SMTP_PORT # default is smtp.163.com is 465
   SMTP_FROM_ADDRESS=YOUR_SMTP_FROM_ADDRESS # example: xxx@163.com
   SMTP_FROM_NAME=gin-auth # example: gin-auth
   
   # enable log
   LOG_ENABLE=false
   ```

4.  Generate graphql related code

   ``` shell
   sh gql.sh
   ```

5. Run Databases

   ```sh
   # minio for windows or you can use docker
   $env:MINIO_ROOT_USER = "admin"
   $env:MINIO_ROOT_PASSWORD = "12345678"
   D:\App\MinIO\minio.exe server D:\App\MinIO\Data --console-address ":9001"
   
   # this is general command for mc.exe, run in cmd
   D:\App\MinIO\mc.exe alias set 'myminio' 'http://ip:port' 'USER' 'PASSWORD'
   
   # redis
   redis-server
   
   # postgresql
   ...
   
   # create database
   CREATE database gin-auth;
   
   # load the two tables
   refer to /resources/*.sql
   
   ```

6. Start the server

   ```sh
   go run main.go
   ```



## Project Structure

This project is organized into several directories and files, each serving a specific purpose. Below is an overview of the project structure:

```
├───.private                # Directory for storing private keys
├───.public                 # Directory for storing public keys
├───main.go                 # Entry point of the application
├───.env                    # Environment variables configuration file
├───go.mod                  # Go module file for dependency management
├───go.sum                  # Go module checksum file
├───gql.sh                  # Script for GraphQL operations
├───gqlgen.yml              # Configuration file for gqlgen
├───.gitignore              # Specifies files and directories to ignore in Git
├───README.md               # Project documentation
├───controllers             # Contains controller logic for handling requests
│   ├───auth                # Authentication-related controllers
│   ├───file                # File handling controllers
│   └───user                # User-related controllers
├───databases               # Database migration and seed files
├───graph                   # GraphQL-related files
│   ├───custom              # Custom GraphQL types and resolvers
│   ├───graphqls            # GraphQL schema definitions
│   ├───model               # GraphQL models
│   └───resolvers           # GraphQL resolvers
├───logs                    # Directory for application logs
├───middlewares             # Middleware functions for request processing
├───models                  # Data models and request structures
│   └───requests            # Request models
├───repositories            # Data access layer for interacting with the database
├───routes                  # Application routing definitions
├───services                # Business logic layer
│   ├───auth                # Authentication services
│   └───user                # User services
└───utils                   # Utility functions and helpers
    ├───consts              # Constants used throughout the application
    ├───cron                # Cron job management
    ├───crypto              # Cryptographic functions
    ├───file                # File handling utilities
    ├───flow                # Workflow management utilities
    ├───jwkmanager          # JSON Web Key management
    ├───jwt                 # JWT (JSON Web Token) handling
    ├───mail                # Email sending utilities
    ├───response            # Response formatting utilities
    └───validation          # Input validation utilities
```

### Overview of Key Directories and Files

- **.private / .public**: These directories store the private and public keys used for signing and verifying JWTs.
- **main.go**: The main entry point of the application where the server is initialized and started.
- **controllers**: Contains the logic for handling incoming requests and returning responses.
- **graph**: Contains all GraphQL-related files, including schema definitions and resolvers.
- **middlewares**: Middleware functions that can be applied to routes for additional processing, such as authentication.
- **models**: Defines the data structures used in the application, including request models.
- **repositories**: Contains the data access layer for interacting with the database.
- **services**: Contains the business logic of the application, separating it from the controller layer.
- **utils**: A collection of utility functions and helpers that are used throughout the application.

This structure promotes a clean separation of concerns, making the codebase easier to navigate and maintain.



## Authentication

Use **refresh token** and **access token**.

- Note that **ACCESS TOKENS ARE NOT STORED IN DATABASE**

- refresh token default 90 days

- access token default 60 minutes

- you can modify in `utils\consts\const.go`, the `JWT_ACCESS_TOKEN_EXPIRY` and `JWT_REFRESH_TOKEN_EXPIRY`

- Refresh tokens are store in `refresh_token` table, you can refer to `resources\refreshToken.sql` for more details.

- Refresh token is random string

- Access token claims, you can modify as you need:

  ```go
  # utils\jwt\jwt.go
  
  publicClaims := jwt.Claims{
      Issuer:  consts.JWT_ISSUER,
      Subject: user.ID,
      // Audience:
      IssuedAt: jwt.NewNumericDate(issuedAt),
      Expiry:   jwt.NewNumericDate(issuedAt.Add(time.Duration(consts.JWT_ACCESS_TOKEN_EXPIRY) * time.Minute)),
  }
  
  // private claims
  privateClaims := map[string]interface{}{
      "email": user.Email,
      // YOU CAN ADD MORE PRIVATE CLAIMS HERE
  }
  ```

- Access tokens are generated based on JWK (JSON Web Key) and are signed to ensure their integrity and authenticity.

- The access token contains public claims such as the issuer, subject (user ID), issued at time, and expiry time. The expiry time is set based on the `JWT_ACCESS_TOKEN_EXPIRY` constant, which you can adjust as needed.

- The private claims can include additional user-specific information, such as the user's email, and can be extended to include other relevant data.

- Access tokens are typically sent in the `Authorization` header of HTTP requests as a Bearer token:

  ```
  Authorization: Bearer <access_token>
  ```

- Ensure to handle token expiration properly in your application. When an access token expires, the client should use the refresh token to obtain a new access token.

- Access tokens are stateless and do not require server-side storage, making them suitable for distributed systems.

- To **decrypt** and **validate** the token, you can use the `ParseToken` and `ParseJWTClaims` functions in `utils/jwt/jwt.go`. These functions will extract the claims from the token and verify its signature using the appropriate public key.

- The keys used for signing the tokens are automatically updated on a scheduled basis. The key rotation is handled by a cron job that runs every Sunday, ensuring that your application uses fresh keys for signing tokens. This enhances security by limiting the lifespan of any given key.

- For more details on how to implement token generation, validation, and key management, refer to the `utils/jwt/jwt.go` and `utils/jwkmanager/jwk.go` files.



## Request Parameter Validation

To ensure that incoming requests contain valid data, this project implements a robust validation mechanism for request parameters. The validation is handled using the `go-playground/validator` package, which allows for flexible and customizable validation rules.

### Validation Process

1. **Request Structs**:
   - Each request type is defined as a struct in the `models/requests` package. These structs include validation tags that specify the required fields and validation rules. For example:

   ```go
   type UserRegisterRequest struct {
       Username string `json:"username" form:"username" validate:"required"`
       Email    string `json:"email" form:"email" validate:"required,email"`
       Password string `json:"password" form:"password" validate:"required,min=6"`
   }
   ```

2. **Global Validator**:
   - A global validator instance is created in `models/requests/validator.go`, which is used to validate the request structs. The `FormatError` function formats validation errors and provides custom error messages based on the validation tags.

3. **Validation Methods**:
   - Each request struct implements a `Validate` method that calls the global validator and formats any errors. This method can be called in both REST API and GraphQL resolvers to ensure that the incoming data is valid.

4. **Binding and Validating**:
   - The `BindAndValidate` function in `utils/validation/validation.go` is used to bind incoming request data to the request struct and validate it. This function works for both REST API and GraphQL requests, ensuring a consistent validation approach across the application.

   ```go
   func BindAndValidate(c *gin.Context, request interface{}) error {
       if err := c.ShouldBind(request); err != nil {
           response.BadRequestWithMessage(c, err.Error())
           return err
       }
       if err := request.(interface{ Validate() error }).Validate(); err != nil {
           response.BadRequestWithMessage(c, err.Error())
           return err
       }
       return nil
   }
   ```

5. **Custom Error Messages**:
   - Custom error messages for validation failures are defined in the `models/requests/auth.go` file. This allows for user-friendly error messages that can be returned in the response when validation fails.

### Example Usage

In the REST API controller, the validation is performed as follows:

```go
func Register(c *gin.Context) {
    var request requests.UserRegisterRequest

    if err := validation.BindAndValidate(c, &request); err != nil {
        return
    }

    err := authService.Register(&request)
    if err != nil {
        response.BadRequestWithMessage(c, err.Error())
        return
    }

    response.Success(c)
}
```

In the GraphQL resolver, the validation is similarly invoked:

```go
func (r *mutationResolver) Register(ctx context.Context, input requests.UserRegisterRequest) (bool, error) {
    if err := input.Validate(); err != nil {
        return false, err
    }

    err := authService.Register(&input)
    if err != nil {
        return false, err
    }

    return true, nil
}
```

By implementing this validation mechanism, the application ensures that all incoming requests are properly validated, reducing the risk of errors and improving overall reliability.



## References

go gorm

https://gorm.io/docs/query.html

Github Action + Docker Hub
https://www.gclhaha.top/building/dockerhub.html#%E5%89%8D%E6%8F%90%E6%9D%A1%E4%BB%B6

argon2

 https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

unix timestamp converter

https://www.unixtimestamp.com/


## DB Migrate


在大型项目中，使用 迁移工具 管理 MongoDB 的索引，可以确保数据库结构在不同环境（如开发、测试、生产）中的一致性。以下是具体步骤，展示如何使用 migrate 工具（比如 golang-migrate）实现 MongoDB 索引的迁移管理。

1. 安装 golang-migrate 工具
golang-migrate 是一个常用的数据库迁移工具，支持多种数据库，包括 MongoDB。

安装 CLI

``` sh
go install -tags 'mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

安装 Go 依赖
在你的 Go 项目中安装 golang-migrate 的 Go 模块：

```bash
go get -u github.com/golang-migrate/migrate/v4
```

2. 配置迁移文件夹
创建迁移文件夹
在项目根目录下创建一个专门用于存放迁移脚本的文件夹：

```bash
mkdir db/migrations
```

3. 创建迁移文件
使用 CLI 创建迁移文件
golang-migrate 提供了命令行工具来创建空的迁移脚本：

```bash
从 C:\Users\user\go\bin copy migrate.exe 至  D:\App\Go1.23.2\bin
migrate create -ext sql -dir db/migrations -seq create_user_indexes

or 

go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir migrations -seq create_user_indexes
```
这会在 db/migrations 文件夹中生成两个文件：

xxx_create_users_indexes.up.sql（向数据库添加索引的逻辑）
xxx_create_users_indexes.down.sql（回滚索引的逻辑）

4. 编写 MongoDB 的迁移脚本
由于 MongoDB 使用的是 JSON 或 BSON 操作而不是 SQL，需要使用 migrate 提供的 MongoDB 驱动（-database mongodb），迁移文件需要用 Go 的逻辑编写。

编写 .up.go（创建索引）
以下是一个示例迁移脚本，添加唯一索引到 users 集合的 email 字段：

```go
package migrations

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Up 添加索引
func Up(db *mongo.Database) error {
	collection := db.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // 对 email 字段创建索引
		Options: options.Index().SetUnique(true),
	}

	// 创建索引
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("Failed to create index on 'email': %v", err)
		return err
	}
	log.Println("Unique index on 'email' created successfully")
	return nil
}
```

编写 .down.go（回滚索引）
以下是一个回滚脚本，删除 users 集合的 email 索引：

```go
package migrations

import (
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
)

// Down 删除索引
func Down(db *mongo.Database) error {
	collection := db.Collection("users")

	// 删除索引
	_, err := collection.Indexes().DropOne(context.Background(), "email_1") // 索引名默认为字段加 "_1"
	if err != nil {
		log.Printf("Failed to drop index on 'email': %v", err)
		return err
	}
	log.Println("Unique index on 'email' dropped successfully")
	return nil
}
```

运行以下命令，将数据库的迁移版本初始化：

```bash
migrate -database mongodb://localhost:27017/my_database -path db/migrations up
```
-database: MongoDB 数据库连接字符串。
-path: 迁移文件的路径。
回滚迁移
如果需要回滚最近的迁移（例如删除索引），运行：

```bash
migrate -database mongodb://localhost:27017/my_database -path db/migrations down
```
检查迁移状态
运行以下命令查看当前的迁移状态：

```bash
migrate -database mongodb://localhost:27017/my_database -path db/migrations version
```

6. 在 Gin 项目中集成
初始化迁移脚本
在服务启动时，自动执行数据库迁移：

```go
package main

import (
	"log"
	"myapp/database"
	"myapp/migrations"
)

func main() {
	// 初始化 MongoDB 数据库
	db, err := database.InitMongoDB("mongodb://localhost:27017", "my_database")
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	// 执行迁移
	if err := migrations.Up(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 启动 Gin 服务
	r := setupRouter()
	r.Run(":8080")
}
```

7. 推荐的最佳实践
版本控制：

确保迁移文件版本号一致，避免不同环境中的索引结构不一致。
使用 Git 管理迁移文件，便于团队协作。
环境隔离：

在开发、测试和生产环境中使用不同的数据库连接字符串，防止误操作。
自动化部署：

将迁移步骤集成到 CI/CD 流程中，确保部署时自动应用最新迁移。
回滚测试：

在开发环境中定期测试迁移的回滚逻辑，确保万一需要回滚时能够正常执行。



claude

migrate CLI

安装

你可以使用 migrate CLI 工具来创建迁移文件。
首先安装 migrate CLI 工具：
``` bash
# 使用 go install
go install -tags 'mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 或者在 MacOS 上使用 brew
brew install golang-migrate
```
然后使用以下命令创建迁移文件：
```bash
migrate create -ext json -dir migrations -seq create_user_indexes
```
这个命令会自动：
这个命令会自动：
创建 migrations 目录（如果不存在）
生成两个文件：
000001_create_user_indexes.up.json
000001_create_user_indexes.down.json
然后你只需要编辑这两个文件，添加相应的 MongoDB 命令即可。
up.json 是执行迁移时的操作（创建索引）
down.json 是回滚时的操作（删除索引）
这样可以确保迁移文件的命名规范和版本号的正确性。