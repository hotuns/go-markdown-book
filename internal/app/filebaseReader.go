package app

import (
	"log"
	"os"
	"time"

	"github.com/hedongshu/go-md-book/internal/types"
	"github.com/hedongshu/go-md-book/internal/utils"
	"github.com/kataras/iris/v12"
)

func HomeHandler(ctx iris.Context) {
	// 记录耗时
	defer utils.TimeTrack(time.Now(), "homeHandler")

	if err := ctx.View("home.html"); err != nil {
		log.Println(err)
	}
}

func CategoriesHandler(ctx iris.Context) {
	// 记录耗时
	defer utils.TimeTrack(time.Now(), "categoriesHandler")

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

func ArticleHandler(ctx iris.Context) {
	// 记录耗时
	defer utils.TimeTrack(time.Now(), "articleHandler")

	f := ctx.Params().Get("f")

	// 设置 Gitalk ID
	Gitalk.Id = utils.MD5(f)
	ctx.ViewData("Gitalk", Gitalk)
	ctx.ViewData("ActiveNav", f)

	if utils.IsInSlice(IgnoreFile, f) {
		return
	}

	mdfile := MdDir + "/" + f + ".md"

	_, err := os.Stat(mdfile)
	if err != nil {
		ctx.StatusCode(404)
		ctx.Application().Logger().Errorf("Not Found '%s', Path is %s", mdfile, ctx.Path())
		return
	}

	bytes, err := os.ReadFile(mdfile)
	if err != nil {
		ctx.StatusCode(500)
		ctx.Application().Logger().Errorf("ReadFile Error '%s', Path is %s", mdfile, ctx.Path())
		return
	}

	ctx.ViewData("Article", utils.MdToHtml(bytes, TocPrefix))

	ctx.View("article.html")
}

func GetAllMarkDownsFromFile() {
	var option utils.Option
	option.RootPath = []string{MdDir}
	option.SubFlag = true
	option.IgnorePath = IgnorePath
	option.IgnoreFile = IgnoreFile
	tree, _ := utils.Explorer(option)

	TreeArticles := make([]types.TreeArticle, 0)
	Categories := make([]string, 0)
	Articles := make([]types.Article, 0)

	for _, v := range tree.Children {
		for _, item := range v.Children {
			if item.IsDir {
				Categories = append(Categories, item.Name)
			}

			thelist := getArticles(item)
			Articles = append(Articles, thelist...)

			TreeArticles = append(TreeArticles, types.TreeArticle{
				CategorieName: item.Name,
				List:          thelist,
			})
		}
	}

	GlobleDatas.Articles = Articles
	GlobleDatas.Categories = Categories
	GlobleDatas.TreeArticles = TreeArticles
}

func getArticles(node *utils.Node) []types.Article {
	list := make([]types.Article, 0)

	if !node.IsDir {
		info := utils.GetArticleInfo(*node)
		list = append(list, info)
	}

	if len(node.Children) > 0 {
		for _, v := range node.Children {
			newinfo := getArticles(v)
			list = append(list, newinfo...)
		}
	}

	return list
}
