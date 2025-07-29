# golang

# 1. Components
## 1.1 User Service (Graphql) (Gorm + Postgresql)

### 1.1.1 Introduction about graphql
Graphql in go use gqlgen
```
go run github.com/99designs/gqlgen init
```

#### Description from gqlgen about generated files
- `gqlgen.yml` - The gqlgen config file, knobs for controlling the generated code.
- `graph/generated/generated.go` - The GraphQL execution runtime, the bulk of the generated code.
- `graph/model/models_gen.go` - Generated models required to build the graph. Often you will override these with your own models. Still very useful for input types.
- `graph/schema.graphqls` - This is the file where you will add GraphQL schemas.
- `graph/schema.resolvers.go` - This is where your application code lives. generated.go will call into this to get the data the user has requested.
---

#### Key components
- ***Query***:
    - a query in graphql is asking for data, you use a query and specify what you want and graphql will return it back to you.
    - Example: 

```bash
// schema.resolvers.go
func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
  var links []*model.Link
  dummyLink := model.Link{
    Title: "our dummy link",
    Address: "https://address.org",
    User: &model.User{Name: "admin"},
  }
	links = append(links, &dummyLink)
	return links, nil
}

// query 
query {
	links{
    title
    address,
    user{
      name
    }
  }
}

// value get
{
  "data": {
    "links": [
      {
        "title": "our dummy link",
        "address": "https://address.org",
        "user": {
          "name": "admin"
        }
      }
    ]
  }
}
```
> Notice that this function takes a Context and returns slice of Links and an error(if there are any). ctx argument contains the data from the person who sends request like which user is working with app(we’ll see how later), etc.


- ***Mutation***
    - Are just like queries but they can cause a data write, Technically Queries can be used to write data too however it’s not suggested to use it. So mutations are like queries, they have names, parameters and they can return data.

    - Example
```bash
// schema.resolvers.go
func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	var link model.Link
	var user model.User
	link.Address = input.Address
	link.Title = input.Title
	user.Name = "test"
	link.User = &user
	return &link, nil
}

// mutation
mutation {
  createLink(input: {title: "new link", address:"http://address.org"}){
    title,
    user{
      name
    }
    address
  }
}


// data get
{
  "data": {
    "createLink": {
      "title": "new link",
      "user": {
        "name": "test"
      },
      "address": "http://address.org"
    }
  }
}
```

- ***Convert value from graphql - model and again***
    - Example

```bash
// internal/user/users.go
package users

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

// internal/links/links.go
package links

import (
	database "github.com/glyphack/go-graphql-hackernews/internal/pkg/db/mysql"
	"github.com/glyphack/go-graphql-hackernews/internal/users"
	"log"
)

// #1
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

//#2
func (link Link) Save() int64 {
	//#3
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.Exec(link.Title, link.Address)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted!")
	return id
}


// schema.resolvers.go
func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	var link links.Link
	link.Title = input.Title
	link.Address = input.Address
	linkID := link.Save()
	return &model.Link{ID: strconv.FormatInt(linkID, 10), Title:link.Title, Address:link.Address}, nil
}

// query
mutation create{
  createLink(input: {title: "something", address: "somewhere"}){
    title,
    address,
    id,
  }
}

// data get
{
  "data": {
    "createLink": {
      "title": "something",
      "address": "somewhere",
      "id": "1"
    }
  }
}
```

- ***Authenticaiton***
    - Use *github.com/dgrijalva/jwt-go* library
    - Example

```bash
// pkg/jwt/jwt.go
package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// secret key being used to sign tokens
var (
	SecretKey = []byte("secret")
)



// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}



// internal/users/users.go (sign up and login)
package users

import (
	"database/sql"
	"github.com/glyphack/go-graphql-hackernews/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"

	"log"
)

type User struct {
	ID       string `json:"id"`
	Username     string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create() {
	statement, err := database.Db.Prepare("INSERT INTO Users(Username,Password) VALUES(?,?)")
	print(statement)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}



// HashPassword hashes given password 
(a request comes to our resolver, we need to know which user is sending the request. To accomplish this, we have to write middleware that’s executed before the request reaches the resolver. This middleware resolves the user from the incoming request and passes this on to the resolver.)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


// internal/users/users.go (authen middleware )
//GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}



// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/glyphack/go-graphql-hackernews/internal/users"
	"github.com/glyphack/go-graphql-hackernews/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// create user and check if user exists in db
			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.ID = strconv.Itoa(id)
			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}



// server/server.go
package main

import (
	"github.com/glyphack/go-graphql-hackernews/internal/auth"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	hackernews "github.com/glyphack/go-graphql-hackernews"
	"github.com/glyphack/go-graphql-hackernews/internal/pkg/db/mysql"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	database.InitDB()
	database.Migrate()
	server := handler.NewDefaultServer(hackernews.NewExecutableSchema(hackernews.Config{Resolvers: &hackernews.Resolver{}}))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
```

## 1.2 Team Service (Rest API) (GIN + GORM + Postgresql)

---
## 1.3 Asset Service (Rest API) (GIN + GORM + Postgresql)

