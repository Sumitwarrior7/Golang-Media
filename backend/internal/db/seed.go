package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"math/rand"

	"github.com/Sumitwarrior7/social/internal/store"
)

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	log.Println("Seeding started!!!!")

	// Users are getting created via transactions i.e., atomic in nature
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating User :", err)
			return
		}
	}
	tx.Commit()
	log.Println("Added Users!!!!")

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating Post :", err)
			return
		}
	}
	log.Println("Added Posts!!!!")

	comments := generatComments(400, posts, users)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating Comment :", err)
			return
		}
	}
	log.Println("Added Comments!!!!")
	log.Println("Seeding Completed!!!")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserId:  user.Id,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generatComments(num int, posts []*store.Post, users []*store.User) []*store.Comment {
	randomComments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		randomPost := posts[rand.Intn(len(posts))]
		randomUser := users[rand.Intn(len(users))]
		randomComments[i] = &store.Comment{
			PostId:  randomPost.Id,
			UserId:  randomPost.UserId,
			Content: comments[rand.Intn(len(comments))],
			User:    *randomUser,
		}

	}
	return randomComments
}
