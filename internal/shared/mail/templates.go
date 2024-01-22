package mail

import "fmt"

func (m *MailService) createActivationTemplate(link string, email string) string {
	l := m.config.AppLink + link
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
				<p style="font-weight: 600">Спасибо за регистрацию в магазине одежды Family Moda!</p>
				<a href="%s">подтвердить регистрацию</a>
			</div>
		</div>
	</body>
</html>
    `, email, l)
}

func (m *MailService) createChangePasswordCodeTemplate(code string, email string) string {
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
			<h1>Подтверждение смены пароля</h1>
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
				<strong>%s</strong>
			</div>
		</div>
	</body>
</html>
    `, email, code)
}

func (m *MailService) createOrderActivationTemplate(link string, email string) string {
	l := m.config.AppLink + link
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
			<h1>Подтверждение оформления заказа</h1>
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
				<a href="%s">подтвердить заказ</a>
			</div>
		</div>
	</body>
</html>
    `, email, l)
}
