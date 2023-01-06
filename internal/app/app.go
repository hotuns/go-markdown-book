package app

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hedongshu/go-md-book/internal/api"
	"github.com/hedongshu/go-md-book/internal/bindata/assets"
	"github.com/hedongshu/go-md-book/internal/bindata/views"
	"github.com/hedongshu/go-md-book/internal/types"
	"github.com/hedongshu/go-md-book/internal/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"github.com/urfave/cli/v2"
)

var (
	MdDir      string
	Env        string
	Title      string
	Title2     string
	Index      string
	LayoutFile = "layouts/layout.html"
	LogsDir    = "cache/logs/"
	TocPrefix  = "[toc]"
	IgnoreFile = []string{`favicon.ico`, `.DS_Store`, `.gitignore`, `README.md`}
	IgnorePath = []string{`.git`}
	Cache      time.Duration
	Analyzer   types.Analyzer
	Gitalk     types.Gitalk
)

// web服务器默认端口
const DefaultPort = 5006

func RunWeb(ctx *cli.Context) error {
	initParams(ctx)

	app := iris.New()

	setLog(app)

	tmpl := iris.HTML(views.AssetFile(), ".html").Reload(true)

	app.RegisterView(tmpl)
	tmpl.AddFunc("getChildrenCount", func(node utils.Node) int {
		return len(node.Children)
	})
	app.OnErrorCode(iris.StatusNotFound, api.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, api.InternalServerError)

	setIndexAuto := false
	if Index == "" {
		setIndexAuto = true
	}

	app.Use(func(ctx iris.Context) {
		activeNav := getActiveNav(ctx)

		Categories, Articles, Tree_articles := getPosts(activeNav)

		if setIndexAuto {
			Index = "/"
		}

		// 设置 Gitalk ID
		Gitalk.Id = utils.MD5(activeNav)

		ctx.ViewData("Gitalk", Gitalk)
		ctx.ViewData("Analyzer", Analyzer)
		ctx.ViewData("Title", Title)
		ctx.ViewData("Title2", Title2)
		ctx.ViewData("Articles", Articles)
		ctx.ViewData("TreeArticles", Tree_articles)
		ctx.ViewData("ActiveNav", activeNav)
		ctx.ViewData("Categories", Categories)
		ctx.ViewLayout(LayoutFile)

		ctx.Next()
	})

	app.Favicon("./favicon.ico")
	app.HandleDir("/static", assets.AssetFile())

	app.Get("/", iris.Cache(Cache), homeHandler)
	app.Get("/article/{f:path}", iris.Cache(Cache), articleHandler)
	app.Get("/categories", iris.Cache(Cache), categoriesHandler)
	app.Get("/categories/{f:path}", iris.Cache(Cache), categoriesHandler)

	app.Run(iris.Addr(":" + strconv.Itoa(parsePort(ctx))))

	return nil
}

func initParams(ctx *cli.Context) {
	MdDir = ctx.String("dir")
	if strings.TrimSpace(MdDir) == "" {
		log.Panic("Markdown files folder cannot be empty")
	}
	MdDir, _ = filepath.Abs(MdDir)

	Env = ctx.String("env")
	Title = ctx.String("title")
	Title2 = ctx.String("title2")
	Index = ctx.String("index")

	Cache = time.Minute * 0
	if Env == "prod" {
		Cache = time.Minute * 3
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

func getPosts(activeNav string) ([]string, []utils.Node, []utils.Node) {
	var option utils.Option
	option.RootPath = []string{MdDir}
	option.SubFlag = true
	option.IgnorePath = IgnorePath
	option.IgnoreFile = IgnoreFile
	tree, _ := utils.Explorer(option)

	Tree_articles := make([]utils.Node, 0)
	Categories := make([]string, 0)
	Articles := make([]utils.Node, 0)

	for _, v := range tree.Children {
		// // v.Children按照ModTime倒序
		// sort.Slice(v.Children, func(i, j int) bool {
		// 	return v.Children[i].ModTime > v.Children[j].ModTime
		// })

		for _, item := range v.Children {
			if item.IsDir {
				Categories = append(Categories, item.Name)
			}

			appendArticles(item, &Articles, activeNav)
			Tree_articles = append(Tree_articles, *item)
		}
	}

	return Categories, Articles, Tree_articles
}

func appendArticles(node *utils.Node, Articles *[]utils.Node, activeNav string) {
	if !node.IsDir {
		*Articles = append(*Articles, *node)
		if node.Link == "/"+activeNav {
			node.Active = "active"
		}
	}

	if len(node.Children) > 0 {
		for _, v := range node.Children {
			appendArticles(v, Articles, activeNav)
		}
	}
}

func getActiveNav(ctx iris.Context) string {
	f := ctx.Params().Get("f")
	if f == "" {
		f = Index
	}
	return f
}

func homeHandler(ctx iris.Context) {
	activeNav := getActiveNav(ctx)

	List := make([]map[string]interface{}, 0)
	_, Articles, _ := getPosts(activeNav)
	for _, v := range Articles {
		info := utils.GetArticleInfo(v)
		List = append(List, structs.Map(info))
	}

	ctx.ViewData("List", List)
	if err := ctx.View("home.html"); err != nil {
		log.Println(err)
	}
}

func categoriesHandler(ctx iris.Context) {
	categorie := ctx.Params().Get("f")
	showAll := false
	List := make([]utils.Article, 0)

	if categorie == "" {
		showAll = true
	} else {
		_, _, TreeArticles := getPosts(categorie)
		for _, v := range TreeArticles {
			if v.Name == categorie {
				for _, item := range v.Children {
					info := utils.GetArticleInfo(*item)
					List = append(List, info)
				}
				break
			}
		}
	}

	ctx.ViewData("ShowAll", showAll)
	ctx.ViewData("Categorie", categorie)
	ctx.ViewData("List", List)

	ctx.View("categories.html")
}

func articleHandler(ctx iris.Context) {
	f := getActiveNav(ctx)

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

	ctx.ViewData("Article", mdToHtml(bytes))

	ctx.View("article.html")
}

func mdToHtml(content []byte) template.HTML {
	strs := string(content)

	var htmlFlags blackfriday.HTMLFlags

	if strings.HasPrefix(strs, TocPrefix) {
		htmlFlags |= blackfriday.TOC
		strs = strings.Replace(strs, TocPrefix, "<br/><br/>", 1)
	}

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: htmlFlags,
	})

	unsafe := blackfriday.Run([]byte(strs), blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(blackfriday.CommonExtensions))
	html := bluemonday.UGCPolicy().AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code").SanitizeBytes(unsafe)

	return template.HTML(string(html))
}
