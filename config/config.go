package config

// Config is the configuration of the blog-generator
type Config struct {
	Generator struct {
		Repo   string
		Tmp    string
		Branch string
		Dest   string
		UseRSS bool
	}
	Blog struct {
		URL            string
		Language       string
		Description    string
		Dateformat     string
		Title          string
		Author         string
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
