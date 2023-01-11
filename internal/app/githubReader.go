package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/hedongshu/go-md-book/internal/utils"
	"github.com/kataras/iris/v12"
)

var (
	client *github.Client
)

func GetAllMarkDownsFromGithub() {
	fmt.Println("GetAllMarkDownsFromGithub start")

	ctx := context.Background()
	client = github.NewClient(nil)

	opt := github.RepositoryContentGetOptions{
		Ref: "main",
	}
	// 获取分类
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, GithubStr.Owner, GithubStr.Repo, "/", &opt)
	fmt.Println(directoryContent)

	if err == nil {
		fmt.Println(err)
		return
	}

	// TreeArticles := make([]utils.Node, 0)
	Categories := make([]string, 0)
	Articles := make([]utils.Node, 0)

	for _, item := range directoryContent {
		if item.GetType() != "dir" || strings.HasPrefix(item.GetName(), "_") {
			return
		}
		Categories = append(Categories, *item.Name)
	}
	appendGithubArticles(Categories, &Articles)
}

func appendGithubArticles(categories []string, articles *[]utils.Node) {
	ctx := context.Background()
	opt := github.RepositoryContentGetOptions{
		Ref: "main",
	}

	for _, item := range categories {

		// 获取分类下的文章
		_, directoryContent, _, err := client.Repositories.GetContents(ctx, GithubStr.Owner, GithubStr.Repo, "/"+item, &opt)
		if err == nil {
			log.Println(err)
			return
		}
		for _, item2 := range directoryContent {
			if item2.GetType() == "file" {
				*articles = append(*articles, utils.Node{
					Name:     item2.GetName(),
					ShowName: item2.GetName()[:strings.LastIndex(item2.GetName(), ".")],
					IsDir:    false,
					Link:     item2.GetDownloadURL(),
				})
			}
		}

	}
}

func Github_HomeHandler(ctx iris.Context) {

	if err := ctx.View("home.html"); err != nil {
		log.Println(err)
	}
}

func Github_ArticleHandler(ctx iris.Context) {

}

func Github_CategoriesHandler(ctx iris.Context) {

}
