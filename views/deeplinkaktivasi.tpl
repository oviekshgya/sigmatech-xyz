<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>Redirect url in Javascript</title>
</head>
<body>
<p id="kode">{{.Link}}</p>
<script>
    var kode = document.getElementById("kode").innerHTML
    window.location = kode;
</script>
</body>
</html>