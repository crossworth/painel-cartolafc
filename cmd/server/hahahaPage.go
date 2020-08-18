package main

const hahahaPage = `<!doctype html>
<html lang="pt-BR">
<head>
  <meta charset="UTF-8">
  <title>Hahaha</title>
	<link rel="shortcut icon" href="/favicon.ico" type="image/x-icon">
	<link rel="icon" href="/favicon.ico" type="image/x-icon">
	<meta name="robots" content="noindex,nofollow">
</head>
<body>
<br>
<center>
  <img src="/hahaha.gif" alt="hahaha"/>
</center>
<script>
  var a = new Audio('magicword.mp3');
  a.addEventListener('ended', function() {
    this.currentTime = 0;
    this.play();
  }, false);
  a.play();
  document.addEventListener('click', function() {
    a.play();
  });
  var e = '{{ .Message }}';
  if (e.length > 0) {
    alert(e);
  }
</script>
</body>
</html>`
