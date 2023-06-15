package main
import(
  "fmt"
  "math/rand"
  "github.com/gin-gonic/gin"
  "github.com/go-redis/redis"
  // "strings"
  "time"
)

type Input struct{
  Url string `json:"url"`
}

func init(){
  rand.Seed(time.Now().UnixNano())
}

func main(){
  // Redis Config
  client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
  })
  pong, err := client.Ping().Result()
  if err!=nil{
    fmt.Println("Error connecting to Redis:", err)
    return
  }
  fmt.Println("Connected to Redis:", pong)
  r := gin.Default()
  r.GET("/status", func(c *gin.Context){
    c.JSON(200, gin.H{"status":"Running",})
  })
  r.POST("/shorten", func(c *gin.Context){
    var i Input;
    err := c.BindJSON(&i);
    if err!=nil{
      panic(err)
    }
    fmt.Println(i.Url)
    exists, err := client.Exists(i.Url).Result()
    if err != nil{
      panic(err)
    }
    if exists!=1{
      genString := fmt.Sprint(rand.Int63n(1000))
      err := client.Set(genString, i.Url, 0).Err()
      if err != nil{
        panic(err)
      }else{
        fmt.Println("SET key-value pair!")
      }
      shortenedUrl := fmt.Sprintf("http://localhost:9000/kadme/%s", genString)
      c.String(200, shortenedUrl)
    }else{
      c.String(200, "URL has already been shortened.")
    }
  })
  r.GET("/kadme/:id", func(c *gin.Context){
    id := c.Param("id")
    url, err := client.Get(id).Result()
    fmt.Println(url)
    if err!=nil{
      panic(err)
    }
    c.Redirect(308, url)
  })
  r.Run(":9000")
}
