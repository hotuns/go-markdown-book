package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	"github.com/hedongshu/go-md-book/internal/types"
	"github.com/hedongshu/go-md-book/internal/utils"
	"github.com/kataras/iris/v12"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
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

	if err != nil {
		log.Println(err)
		return
	}

	TreeArticles := make([]types.TreeArticle, 0)
	Categories := make([]string, 0)
	Articles := make([]types.Article, 0)

	for _, item := range directoryContent {
		if strings.HasPrefix(item.GetName(), "_") {
			continue
		}

		if item.GetType() == "dir" {
			Categories = append(Categories, item.GetName())

			theList := getArticlesFromGithub(item, item.GetName())
			Articles = append(Articles, theList...)

			TreeArticles = append(TreeArticles, types.TreeArticle{
				CategorieName: item.GetName(),
				List:          theList,
			})
		}
	}

	GlobleDatas.Articles = Articles
	GlobleDatas.Categories = Categories
	GlobleDatas.TreeArticles = TreeArticles
}

func getArticlesFromGithub(content *github.RepositoryContent, category string) []types.Article {
	list := make([]types.Article, 0)

	if content.GetType() == "file" {
		title := content.GetName()
		i := strings.LastIndex(title, ".")
		if i != -1 {
			title = title[:i]
		}

		info := types.Article{
			Title:       title,
			Category:    category,
			PublishTime: "",
			Link:        "/" + category + "/" + title,
		}
		downUrl := content.GetDownloadURL()
		resp, err := http.Get(downUrl)
		if err != nil {
			log.Fatal(err)
		}
		// 读取返回结果
		body, _ := io.ReadAll(resp.Body)
		unsafe := blackfriday.Run(body)
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
		if err != nil {
			log.Println(err)
		}
		plist := doc.Find("p").Text()
		if strings.Contains(plist, "[toc]") {
			plist = strings.ReplaceAll(plist, "[toc]", "")
		}
		runeList := []rune(plist)
		if len(runeList) > 150 {
			plist = string(runeList[:150])
		}
		info.Preview = plist + "..."

		list = append(list, info)
	}

	if content.GetType() == "dir" {
		ctx := context.Background()
		opt := github.RepositoryContentGetOptions{
			Ref: "main",
		}
		// 获取文章列表
		_, posts, _, err := client.Repositories.GetContents(ctx, GithubStr.Owner, GithubStr.Repo, "/"+content.GetName(), &opt)
		if err != nil {
			return list
		}

		for _, post := range posts {
			newitem := getArticlesFromGithub(post, content.GetName())
			list = append(list, newitem...)
		}

	}

	return list
}

func Github_HomeHandler(ctx iris.Context) {
	if err := ctx.View("home.html"); err != nil {
		log.Println(err)
	}
}

func Github_ArticleHandler(ctx iris.Context) {
	f := ctx.Params().Get("f")
	// 设置 Gitalk ID
	Gitalk.Id = utils.MD5(f)
	ctx.ViewData("Gitalk", Gitalk)
	ctx.ViewData("ActiveNav", f)

	if utils.IsInSlice(IgnoreFile, f) {
		return
	}

	opt := github.RepositoryContentGetOptions{
		Ref: "main",
	}
	// 获取markdwon文件
	fileContent, _, _, err := client.Repositories.GetContents(ctx, GithubStr.Owner, GithubStr.Repo, "/"+f+".md", &opt)

	if err != nil {
		ctx.StatusCode(500)
		ctx.Application().Logger().Errorf("ReadFile Error '%s', Path is %s", f, ctx.Path())
		return
	}

	c, _ := fileContent.GetContent()
	ctx.ViewData("Article", utils.MdToHtml([]byte(c), TocPrefix))

	ctx.View("article.html")

}

func Github_CategoriesHandler(ctx iris.Context) {
	categorie := ctx.Params().Get("f")
	showAll := false

	if categorie == "" {
		showAll = true
	} else {

		for _, item := range GlobleDatas.TreeArticles {
			if item.CategorieName == categorie {
				ctx.ViewData("theList", item.List)
			}
		}
	}

	ctx.ViewData("ShowAll", showAll)
	ctx.ViewData("Categorie", categorie)

	ctx.View("categories.html")
}
