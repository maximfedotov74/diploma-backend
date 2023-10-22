package mail

import "fmt"

func (m *MailService) CreateActivationTemplate(link string, email string) string {
	l := fmt.Sprintf("%s/api/user/api/activate/%s", m.appLink, link)
	return fmt.Sprintf(`
  <!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
	</head>
	<body>
		<div>
			<h1>Регистрация на сайте одежды</h1>
      <h2>Здравствуйте, уважаемый %s</h2>
			<div
				style="
					background-color: #8e92fa;
					padding: 15px;
					border-radius: 8px;
					color: #fff;
					font-weight: 600;
				"
			>
				<p style="font-weight: 600">Приветствуем вас, дорогой друг!</p>
				<p style="font-weight: 600">Спасибо за регистрацию на нашем форуме!</p>
				<a href="%s">подтвердить регистрацию</a>
			</div>
		</div>
	</body>
</html>
    `, email, l)
}
