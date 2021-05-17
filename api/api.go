package api

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/itchyny/base58-go"
	"github.com/laupski/url-shortener/etcd"
	"github.com/sirupsen/logrus"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	HTTPPort     string `toml:"api_port"`
	WriteTimeout string `toml:"write_timeout"`
	ReadTimeout  string `toml:"read_timeout"`
}

type ShortenPayload struct {
	Link string `form:"link" binding:"required"`
}

func RunApi(c Config, conn etcd.Connection) {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./ui", true)))
	router.GET("/:shortUrl", func(c *gin.Context) {
		Redirect(c, conn)
	})

	/*router.NoRoute(func(c *gin.Context) {
		fmt.Println(os.Getwd())
		//c.HTML(http.StatusNotFound, "./ui/404.html",gin.H{})
	})*/
	/*router.GET("api/health", func(c *gin.Context) {
		c.String(200, "ok")
	})*/

	router.POST("/api/shorten", func(c *gin.Context) {
		Shorten(c, conn)
	})

	wt, err := time.ParseDuration(c.WriteTimeout)
	if err != nil {
		logrus.Fatal(err)
	}
	rt, err := time.ParseDuration(c.ReadTimeout)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := &http.Server{
		Handler: router,
		Addr:    c.HTTPPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: wt,
		ReadTimeout:  rt,
	}

	fmt.Printf("Now serving on port %v\n", c.HTTPPort)
	log.Fatal(srv.ListenAndServe())
}

func Shorten(c *gin.Context, conn etcd.Connection) {
	fmt.Println("Hit /shorten route")
	url := ShortenPayload{}
	err := c.Bind(&url)
	if err != nil {
		c.String(http.StatusOK, "Error: This is not a valid request")
	} else if url.Link == "" {
		c.String(http.StatusOK, "Error: Empty Request")
	} else if isValidURL(url.Link) == false {
		c.String(http.StatusOK, "Error: Invalid link")
	} else {
		key := generateKey(url.Link)
		err = etcd.PutRedirect(conn, key, url.Link)
		if err != nil {
			logrus.Warn(err)
		}

		link := fmt.Sprintf("http://%v/%v", c.Request.Host, key)
		c.String(http.StatusOK, "<a href =\"%v\">%v</a>", link, link)
	}
}

func Redirect(c *gin.Context, conn etcd.Connection) {
	shortUrl := c.Param("shortUrl")
	redirect, err := etcd.GetRedirect(conn, shortUrl)
	if redirect != "" && err == nil {
		c.Redirect(http.StatusFound, redirect)
	} else {
		c.String(http.StatusNotFound, "Redirect not found")
	}
}

func isValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func generateKey(s string) string {
	algorithm := sha256.New()
	algorithm.Write([]byte(s))
	hash := algorithm.Sum(nil)
	number := new(big.Int).SetBytes(hash).Uint64()
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode([]byte(fmt.Sprintf("%d", number)))
	if err != nil {
		fmt.Println(err)
	}
	return string(encoded)[:8]
}
