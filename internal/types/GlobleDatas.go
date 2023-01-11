package types

// Article 文章
type Article struct {
	// 文章标题
	Title string `json:"title"`
	// 分类
	Category string `json:"category"`
	// 发布时间
	PublishTime string `json:"publishTime"`
	// 文章预览， 截取文章前100个字符
	Preview string `json:"preview"`
	// 链接
	Link string `json:"link"`
}

type TreeArticle struct {
	CategorieName string    `json:"categorieName"`
	List          []Article `json:"list"`
}

type GlobleData struct {
	Categories   []string      `json:"categories"`    // 一级目录，即分类
	Articles     []Article     `json:"articles"`      // 扁平化文章列表
	TreeArticles []TreeArticle `json:"tree_articles"` // 树结构文章列表
}
