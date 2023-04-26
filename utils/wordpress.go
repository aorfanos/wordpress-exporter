package utils

type Wordpress struct {
	posts, categories, tags, pages, comments, media, users, adminUsers, themes, plugins, taxonomies int
	MonitoredWordpress string
	UserAgent string
	Auth WPAuth
}

type WPAuth struct {
	Use bool
	Username string
	Password string
}

func NewWordpress(monitor, ua, authuser, authpass string, useAuth bool) *Wordpress {
	return &Wordpress{
		MonitoredWordpress: monitor,
		UserAgent: ua,
		Auth: WPAuth{
			Use: useAuth,
			Username: authuser,
			Password: authpass,
		}, 
	}
}