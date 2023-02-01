package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hedongshu/go-md-book/internal/api"
	"github.com/hedongshu/go-md-book/internal/bindata/assets"
	"github.com/hedongshu/go-md-book/internal/bindata/views"
	"github.com/hedongshu/go-md-book/internal/types"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/urfave/cli/v2"
)

var (
	Origin     string
	MdDir      string
	Env        string
	Title      string
	Title2     string
	LayoutFile               = "layouts/layout.html"
	LogsDir                  = "cache/logs/"
	TocPrefix                = "[toc]"
	IgnoreFile               = []string{`favicon.ico`, `.DS_Store`, `.gitignore`, `README.md`}
	IgnorePath               = []string{`.git`}
	Cache      time.Duration = 3
	Analyzer   types.Analyzer
	Gitalk     types.Gitalk
	GithubStr  types.GithubStr
)

var GlobleDatas types.GlobleData

// web服务器默认端口
const DefaultPort = 5006

func RunWeb(ctx *cli.Context) error {
	initParams(ctx)

	app := iris.New()

	setLog(app)

	tmpl := iris.HTML(views.AssetFile(), ".html").Reload(true)
	app.RegisterView(tmpl)
	tmpl.AddFunc("getArticlesLen", func(articles []types.Article) int {
		return len(articles)
	})
	tmpl.AddFunc("isHomeNavActive", func(url string, name string) string {
		class := ""

		if strings.TrimSpace(url) == strings.TrimSpace(name) {
			class = "current"
		}

		return class
	})
	app.OnErrorCode(iris.StatusNotFound, api.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, api.InternalServerError)
	app.Favicon("./favicon.ico")
	app.HandleDir("/static", assets.AssetFile())

	if Origin == "file" {
		app.Use(func(ctx iris.Context) {

			CurrentPath := ctx.Path()

			GetAllMarkDownsFromFile()

			ctx.ViewData("CurrentPath", CurrentPath)
			ctx.ViewData("Analyzer", Analyzer)
			ctx.ViewData("Title", Title)
			ctx.ViewData("Title2", Title2)

			ctx.ViewData("Categories", GlobleDatas.Categories)
			ctx.ViewData("Articles", GlobleDatas.Articles)
			ctx.ViewData("TreeArticles", GlobleDatas.TreeArticles)

			ctx.ViewLayout(LayoutFile)

			ctx.Next()
		})

		app.Get("/", iris.Cache(Cache), HomeHandler)
		app.Get("/article/{f:path}", iris.Cache(Cache), ArticleHandler)
		app.Get("/categories", iris.Cache(Cache), CategoriesHandler)
		app.Get("/categories/{f:path}", iris.Cache(Cache), CategoriesHandler)
	}
	if Origin == "github" {
		fmt.Println("origin is github")
		GetAllMarkDownsFromGithub()

		app.Use(func(ctx iris.Context) {

			CurrentPath := ctx.Path()

			ctx.ViewData("CurrentPath", CurrentPath)
			ctx.ViewData("Analyzer", Analyzer)
			ctx.ViewData("Title", Title)
			ctx.ViewData("Title2", Title2)
			ctx.ViewData("Articles", GlobleDatas.Articles)
			ctx.ViewData("Categories", GlobleDatas.Categories)
			ctx.ViewData("TreeArticles", GlobleDatas.TreeArticles)
			ctx.ViewLayout(LayoutFile)

			ctx.Next()
		})

		app.Get("/", iris.Cache(Cache), Github_HomeHandler)
		app.Get("/article/{f:path}", iris.Cache(Cache), Github_ArticleHandler)
		app.Get("/categories", iris.Cache(Cache), Github_CategoriesHandler)
		app.Get("/categories/{f:path}", iris.Cache(Cache), Github_CategoriesHandler)

	}

	app.Get("/update", func(ctx iris.Context) {
		if Origin == "file" {
			ctx.HTML("<h3>当前不是github模式</h3>")
			return
		}
		if Origin == "github" {
			oldLen := len(GlobleDatas.Articles)
			GetAllMarkDownsFromGithub()
			nowLen := len(GlobleDatas.Articles)

			ctx.HTML("<h3>更新完成</h3><p>更新前 %d 篇文章, 更新后 %d 篇文章</p><p><a href= '/' >回到首页</a></p>", oldLen, nowLen)
			return
		}
	})

	app.Run(iris.Addr(":" + strconv.Itoa(parsePort(ctx))))

	return nil
}

func initParams(ctx *cli.Context) {
	// 设置文件来源
	Origin = ctx.String("origin")

	if Origin == "github" {
		GithubStr.Owner = ctx.String("github.owner")
		GithubStr.Repo = ctx.String("github.repo")
	}

	MdDir = ctx.String("dir")
	if strings.TrimSpace(MdDir) == "" {
		log.Panic("Markdown files folder cannot be empty")
	}
	MdDir, _ = filepath.Abs(MdDir)

	Env = ctx.String("env")
	Title = ctx.String("title")
	Title2 = ctx.String("title2")

	_cache := ctx.Int("cache")

	Cache = time.Minute * time.Duration(_cache)
	if Env == "dev" {
		Cache = time.Minute * 0
	}

	// 设置分析器
	Analyzer.SetAnalyzer(ctx.String("analyzer-baidu"), ctx.String("analyzer-google"))

	// 设置Gitalk
	Gitalk.SetGitalk(ctx.String("gitalk.client-id"), ctx.String("gitalk.client-secret"), ctx.String("gitalk.repo"), ctx.String("gitalk.owner"), ctx.StringSlice("gitalk.admin"), ctx.StringSlice("gitalk.labels"))

	// 忽略文件
	IgnoreFile = append(IgnoreFile, ctx.StringSlice("ignore-file")...)
	IgnorePath = append(IgnorePath, ctx.StringSlice("ignore-path")...)
}

func setLog(app *iris.Application) {
	os.MkdirAll(LogsDir, 0777)
	f, _ := os.OpenFile(LogsDir+"access-"+time.Now().Format("20060102")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)

	if Env == "prod" {
		app.Logger().SetOutput(f)
	} else {
		app.Logger().SetLevel("debug")
		app.Logger().Debugf(`Log level set to "debug"`)
	}

	// Close the file on shutdown.
	app.ConfigureHost(func(su *iris.Supervisor) {
		su.RegisterOnShutdown(func() {
			f.Close()
		})
	})

	ac := accesslog.New(f)
	ac.AddOutput(app.Logger().Printer)
	app.UseRouter(ac.Handler)
	app.Logger().Debugf("Using <%s> to log requests", f.Name())
}

func parsePort(ctx *cli.Context) int {
	port := DefaultPort
	if ctx.IsSet("port") {
		port = ctx.Int("port")
	}
	if port <= 0 || port >= 65535 {
		port = DefaultPort
	}

	return port
}
