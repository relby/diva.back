package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/relby/diva.back/internal/app"
	"github.com/relby/diva.back/internal/model"
	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"
)

func main() {
	ctx := context.Background()

	diContainer, err := app.NewDIContainer()
	if err != nil {
		panic(err)
	}

	queries, err := diContainer.Queries(ctx)
	if err != nil {
		panic(err)
	}

	customerRespository, err := diContainer.CustomerRepository(ctx)
	if err != nil {
		panic(err)
	}

	cliApp := &cli.App{
		Name:  "diva-excel",
		Usage: "TODO",
		Commands: []*cli.Command{
			{
				Name:    "import",
				Aliases: []string{"i"},
				Args:    true,
				Usage:   "import customers from excel",
				Action: func(ctx *cli.Context) (err error) {
					filename := ctx.Args().First()
					if filename == "" {
						return errors.New("filename not provided")
					}
					file, err := excelize.OpenFile(filename)
					if err != nil {
						return err
					}
					defer func() {
						err = file.Close()
					}()
					rows, err := file.GetRows("Лист1")
					if err != nil {
						return err
					}

					if err := queries.TruncateCustomers(ctx.Context); err != nil {
						panic(err)
					}

					rows = rows[1:]
					for _, row := range rows {
						fullName := row[0]
						phoneNumber := row[1]
						discount, err := strconv.Atoi(strings.Trim(row[2], " %?"))
						if err != nil {
							panic(err)
						}
						customerFullName, err := model.NewCustomerFullName(fullName)
						if err != nil {
							panic(err)
						}

						customerPhoneNumber, err := model.NewCustomerPhoneNumber(phoneNumber)
						if err != nil {
							panic(err)
						}

						customerDiscount, err := model.NewCustomerDiscount(discount)
						if err != nil {
							panic(err)
						}

						customer, err := model.NewCustomer(customerFullName, customerPhoneNumber, customerDiscount)
						if err != nil {
							panic(err)
						}

						if err := customerRespository.Save(ctx.Context, customer); err != nil {
							panic(err)
						}
					}
					return nil
				},
			},
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "export customers to excel",
				Action: func(ctx *cli.Context) error {
					return errors.New("UNIMPLEMENTED")
				},
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
