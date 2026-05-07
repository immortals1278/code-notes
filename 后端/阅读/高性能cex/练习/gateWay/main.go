/*
做了url 反向代理 gin（仅orderService部分）
没做redis client
*/


import{
	"net/url"
	"net/http/httputil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
}

func main() {
	//关闭程序前将所有日志打印到终端
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	orderServiceURL := os.getenv("ORDER_SERVICE")
	if orderServiceURL == "" {
		orderServiceURL = "http://localhost:8100"
	}

	orderService_url,err := url.Parse(orderServiceURL)
	if err != nil {
		logger.Log.Fatal("解析成指针失败", zap.String("url", orderServiceURL), zap.Error(err))
	}

	orderProxy := newReverseProxy(orderService_url)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(metrics.Middleware("gateway"))







	
	func newReverseProxy(URL *url.URL) *httputil.ReverseProxy {
		proxy := httputil.NewSingleHostReverseProxy(URL)

		//重写director逻辑
		originalDirector := proxy.Director 
		proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
		if req.TLS != nil {
			req.Header.Set("X-Forwarded-Proto", "https")
			return
		}
		req.Header.Set("X-Forwarded-Proto", "http")
		}

		//重写errorHandler逻辑
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Log.Error("gateway: 反向代理失败", zap.String("path", r.URL.Path), zap.Error(err))
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error":   "upstream unavailable",
			"service": "gateway",
		})
		}
		return proxy
	}
}