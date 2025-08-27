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

- ***Convert value from graphql - model and value return***
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

#### Example
```bash
// login
mutation {
  login(input: {
    email: "falgod@example.com"
    password: "mypassword123"
  }) {
    token
    user {
      userID
      username
      email
      role
    }
  }
}


// create user
mutation {
  createUser(input: {
    username: "newuser03",
    email: "newuser03@example.com",
    password: "newuser03",
    role: "member"
  }) {
    userID
    username
    email
    role
  }
}

{
  "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI0ZjhlY2RiMi1iZTVmLTQ1NzctYWE4YS01ZmMzMzFjZTY2OTIiLCJyb2xlIjoibWFuYWdlciIsImV4cCI6MTc1Mzg5MDAzNn0.K28NhNIUG3_RRIF6f-b0-QMcaq46psLhczdpPeR9rWc"
}


// fetch user
query {
  fetchUsers {
    userID
    username
    email
    role
  }
}

```

## 1.2 Team Service (Rest API) (GIN + GORM + Postgresql)

```bash
# 1. Create a team
curl -X POST http://localhost:8081/api/v1/teams \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "teamName": "Development Team",
    "managers": [
      {
        "managerId": "4f8ecdb2-be5f-4577-aa8a-5fc331ce6692",
        "managerName": "Main Manager"
      }
    ],
    "members": [
      {
        "memberId": "member-uuid-1",
        "memberName": "John Doe"
      },
      {
        "memberId": "member-uuid-2",
        "memberName": "Jane Smith"
      }
    ]
  }'

# 2. Get all teams
curl -X GET http://localhost:8081/api/v1/teams \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 3. Get specific team
curl -X GET http://localhost:8081/api/v1/teams/TEAM_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 4. Add member to team
curl -X POST http://localhost:8081/api/v1/teams/TEAM_ID/members \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "memberId": "new-member-uuid",
    "memberName": "New Member"
  }'

# 5. Remove member from team
curl -X DELETE http://localhost:8081/api/v1/teams/TEAM_ID/members/MEMBER_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 6. Add manager to team (only main manager can do this)
curl -X POST http://localhost:8081/api/v1/teams/TEAM_ID/managers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "managerId": "new-manager-uuid",
    "managerName": "New Manager"
  }'

# 7. Remove manager from team (only main manager can do this)
curl -X DELETE http://localhost:8081/api/v1/teams/TEAM_ID/managers/MANAGER_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---
## 1.3 Asset Service (Rest API) (GIN + GORM + Postgresql)


# Grafana

![alt text](image.png)

# Example for kafka and redis
[Watch Demo](./kafka-video.mp4)