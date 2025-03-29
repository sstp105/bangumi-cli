package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
)

var username string
var subjectType int
var collectionType int
var createFolders bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List of user subject collections from bangumi.tv",
	Run: func(cmd *cobra.Command, args []string) {
		listHandler()
	},
}

func init() {
	listCmd.Flags().StringVarP(&username, "username", "u", "", "The username of bangumi.tv")
	listCmd.Flags().IntVarP(&subjectType, "type", "t", 2,
		"1 = 书籍 (book)\n"+
			"2 = 动画 (anime)\n"+
			"3 = 音乐 (music)\n"+
			"4 = 游戏 (game)\n"+
			"6 = 三次元 (real)")
	listCmd.Flags().IntVarP(&collectionType, "status", "s", 1,
		"1 = 想看 (wish)\n"+
			"2 = 看过 (collect)\n"+
			"3 = 在看 (do)\n"+
			"4 = 搁置 (on_hold)\n"+
			"5 = 抛弃 (dropped)")
	listCmd.Flags().BoolVarP(&createFolders, "create-folders", "c", false, "Create a folder for each subject")

	rootCmd.AddCommand(listCmd)
}

func listHandler() {
	if username == "" {
		log.Fatal("error: please provide a username using the --username flag.")
	}

	client := bangumi.NewClient()
	resp, err := client.GetUserCollections(username, subjectType, collectionType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("共%d条结果", len(resp))

	if createFolders {
		for _, subject := range resp {
			folderName := subject.Subject.NameCN
			err := createFolder(folderName)
			if err != nil {
				log.Printf("failed to create folder for %s: %v\n", folderName, err)
			}
		}
	}
}

func createFolder(folderName string) error {
	err := os.Mkdir(folderName, 0755)
	if err != nil {
		return err
	}
	return nil
}

