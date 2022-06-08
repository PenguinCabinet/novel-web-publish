package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	global_init()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "make new project here.",
				Action: func(c *cli.Context) error {
					fmt.Print("Title>")
					Title := input_text()
					new_project(Title)
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "make a episode in the project.",
				Action: func(c *cli.Context) error {
					episodes := load_episodes_in_project()
					if c.Args().Len() == 0 {
						fmt.Println("You need to give the sub-title of the episode as First Args.")
					}
					add_episode_in_project(episodes, c.Args().First())
					return nil
				},
			},
			{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "Deploy the project to each website.",
				Action: func(c *cli.Context) error {
					Project_Setting_data, summary := load_project()
					episodes := load_episodes_in_project()
					for _, e := range Project_Setting_data.Deploys {
						if e == "narou" {
							fmt.Println("Deploying to なろう...")
							deploy_of_Narou(episodes, Project_Setting_data, summary)
							fmt.Println("Successd")
						}
					}

					return nil
				},
			},
			{
				Name:    "narou",
				Aliases: []string{},
				Usage:   "narou subcommands",
				Subcommands: []*cli.Command{
					{
						Name:    "login",
						Aliases: []string{"l"},
						Usage:   "login of Narou",
						Action: func(c *cli.Context) error {
							fmt.Print("Email>")
							id := input_text()
							fmt.Print("Password>")
							password := input_password()
							err := login_narou(id, password)
							if err != nil {
								log.Fatalln(err)
							}
							return nil
						},
					},
					{
						Name:    "novel_list",
						Aliases: []string{"nl"},
						Usage:   "print the list of novels published in narou",
						Action: func(c *cli.Context) error {
							narou_secrets_data := load_narou_secret()
							Ret := Get_list_of_narou(&narou_secrets_data)
							novel_list := my_obj_to_json(Ret)
							fmt.Printf("%s\n", novel_list)
							return nil
						},
					},
					{
						Name:    "episode_list",
						Aliases: []string{"el"},
						Usage:   "print the list of episodes of the novel published in narou",
						Action: func(c *cli.Context) error {
							if c.Args().Len() == 0 {
								fmt.Println("You need to give the ncode of the novel as First Args.")
							}
							narou_secrets_data := load_narou_secret()
							Ret := Get_list_of_episode_of_narou(c.Args().First(), &narou_secrets_data)

							novel_list := my_obj_to_json(Ret)
							fmt.Printf("%s\n", novel_list)
							return nil
						},
					},
				},
			},
		},
	}

	err_CLI := app.Run(os.Args)
	if err_CLI != nil {
		log.Fatal(err_CLI)
	}
}
