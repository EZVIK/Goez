package dto

type LoginParams struct {
	Username string `binding:"required" json:"username" form:"username"`
	Password string `binding:"required" json:"password" form:"password"`
	//Nickname string `binding:"required" json:"nickname" form:"nickname"`
}

type RegisterParams struct {
	Username string `binding:"required" json:"username" form:"username"`
	Password string `binding:"required" json:"password" form:"password"`
	RePassword string `binding:"required" json:"re_password" form:"re_password"`
	Nickname string `binding:"required" json:"nickname" form:"nickname"`
}


type ArticleSearchParams struct {
	Field 			string		`binding:"required" json:"field"`
	Value       	string 		`binding:"required" json:"value"`
}

type AddArticleParams struct {
	Title 		string			`binding:"required" json:"title" form:"title"`
	Desc       	string 			`binding:"required" json:"desc" form:"desc"`
	Content     string 			`binding:"required" json:"content" form:"content"`
	Tags       	[]string 		`binding:"required" json:"tags" form:"tags"`
}


type AddTagsParams struct {
	Name        string			`json:"name" form:"name"`
	Names       []string 		`json:"names" form:"names"`

}