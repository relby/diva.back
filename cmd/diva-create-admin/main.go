package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/relby/diva.back/internal/app"
	"github.com/relby/diva.back/internal/model"
)

var (
	login    = flag.String("login", "", "")
	password = flag.String("password", "", "")
	fullName = flag.String("full-name", "", "")
)

func main() {
	flag.Parse()
	if *login == "" {
		log.Fatalln("provide `login` flag")
	}
	if *password == "" {
		log.Fatalln("provide `password` flag")
	}
	if *fullName == "" {
		log.Fatalln("provide `full-name` flag")
	}

	adminLogin, err := model.NewAdminLogin(*login)
	if err != nil {
		panic(err)
	}

	adminHashedPassword, err := model.NewAdminHashedPasswordFromPassword(*password)
	if err != nil {
		panic(err)
	}

	adminFullName, err := model.NewUserFullName(*fullName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	diContainer, err := app.NewDIContainer()
	if err != nil {
		panic(err)
	}

	adminRepository, err := diContainer.AdminRepository(ctx)
	if err != nil {
		panic(err)
	}

	admin, err := model.NewAdminWithRandomID(adminFullName, adminLogin, adminHashedPassword)
	if err != nil {
		panic(err)
	}

	err = adminRepository.Save(ctx, admin)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfuly created admin: %+v\n", admin)
}
