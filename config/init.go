package config

// Config is the configuration of the blog-generator
type Config struct {
	Generator struct {
		Repo string
		Tmp  string
		Dest string
	}
	Blog struct {
		URL            string
		Language       string
		Description    string
		Dateformat     string
		Title          string
		Frontpageposts int
		Statics        struct {
			Files []struct {
				Src  string
				Dest string
			}
			Templates []struct {
				Src  string
				Dest string
			}
		}
	}
}
