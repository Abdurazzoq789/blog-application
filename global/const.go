package global

const (
	dburi = "mongodb://localhost:27017"
	dbname = "blog-application"
	performance = 100
)

var (
	jwtSecret = []byte("blogSecret")
)
