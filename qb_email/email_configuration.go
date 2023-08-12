package qb_email

type SmtpSettingsAuth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type SmtpSettings struct {
	Host    string            `json:"host"`
	Port    int               `json:"port"`
	Secure  bool              `json:"secure"`
	Auth    *SmtpSettingsAuth `json:"auth"`
	From    string            `json:"from"`
	ReplyTo string            `json:"reply_to"`
}
