package typeahead

import "github.com/gin-gonic/gin"

type Filter struct {
	keyword   string
	latitude  string
	longitude string
}

func (f *Filter) transform(c *gin.Context) {
	keyword, _ := c.GetQuery("query")
	latitude := c.GetHeader("latitude")
	longitude := c.GetHeader("longitude")

	f.keyword = keyword
	f.latitude = latitude
	f.longitude = longitude

}
