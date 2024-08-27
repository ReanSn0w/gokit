package smtp

type Config struct {
	SMTP struct {
		Host     string `long:"host" env:"HOST" description:"smtp host value"`
		Port     string `long:"port" env:"PORT" description:"smtp port value"`
		Name     string `long:"name" env:"NAME" description:"smtp name value"`
		Email    string `long:"email" env:"EMAIL" description:"smtp email value"`
		Login    string `long:"login" env:"LOGIN" description:"smtp login value"`
		Password string `long:"password" env:"PASSWORD" description:"smtp password value"`
	} `group:"mongo" namespace:"mongo" env-namespace:"MONGO"`
}

// New - возвращает новый SMTP клиент
func (c Config) NewSMTP() *SMTP {
	return New(
		c.SMTP.Host,
		c.SMTP.Port,
		c.SMTP.Name,
		c.SMTP.Email,
		c.SMTP.Login,
		c.SMTP.Password)
}
