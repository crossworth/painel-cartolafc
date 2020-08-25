package auth

const loginPage = `<!doctype html>
<html lang="pt-BR">
<head>
  <meta charset="UTF-8">
  <title>{{ .Title }}</title>
  <link rel="shortcut icon" href="/favicon.ico" type="image/x-icon">
  <link rel="icon" href="/favicon.ico" type="image/x-icon">
  <meta name="robots" content="noindex,nofollow">
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
</head>
<body>
<br>
<center>
  <a href="/login">Fazer login</a>
</center>
<script>
  var urlParams = new URLSearchParams(window.location.search)
  var reason = urlParams.get('motivo-redirect')
  if (reason !== null && reason !== '') {
    alert(reason)
  }
</script>
</body>
</html>`
