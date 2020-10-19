<!DOCTYPE html>

<html>
<head>
  <title>MicEngine</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <link rel="shortcut icon" href="/static/img/micloud.ico">
  <style type="text/css">
    *,body {
      margin: 0px;
      padding: 0px;
    }

    body {
      margin: 0px;
      font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
      font-size: 14px;
      line-height: 20px;
      background-color: #fff;
    }

    header,
    footer {
      width: 960px;
      margin-left: auto;
      margin-right: auto;
    }

    .logo {
      background-image: url("/static/img/micloud.svg");
      background-repeat: no-repeat;
      -webkit-background-size: 100px 100px;
      background-size: 100px 100px;
      background-position: center center;
      text-align: center;
      font-size: 42px;
      padding: 250px 0 70px;
      font-weight: normal;
      text-shadow: 0px 1px 2px #ddd;
    }

    header {
      padding: 100px 0;
    }

    footer {
      line-height: 1.8;
      text-align: center;
      padding: 50px 0;
      color: #999;
    }

    .description {
      text-align: center;
      font-size: 16px;
    }
    .version {
      text-align: center;
      padding: 0px 0 30px;
    }

    a {
      color: #444;
      text-decoration: none;
    }

    .backdrop {
      position: absolute;
      width: 100%;
      height: 100%;
      box-shadow: inset 0px 0px 100px #ddd;
      z-index: -1;
      top: 0px;
      left: 0px;
    }
  </style>
</head>

<body>
  <header>
    <h1 class="logo">{{.HeaderTitle}}</h1>
    <h2 class="version">版本号：<a href="http://mining-icloud.com/showdoc/web/#/9?page_id=81">{{.Version}}</a></h2>
    <div class="description">
      {{.Copyright}}
    </div>
  </header>
  <footer>
    <div class="author">
      Official website:
      <a href="http://{{.Website}}">{{.Website}}</a> /
      Contact me:
      <a class="email" href="mailto:{{.Email}}">{{.Email}}</a>
    </div>
    <div class="api">
      <a href="http://mining-icloud.com/showdoc/web/#/9?page_id=50">WebAPI接口文档</a> /
      <a href="http://mining-icloud.com/showdoc/web/#/1">Grafana读取平台数据接口文档</a>
    </div>
  </footer>
  <div class="backdrop"></div>
</body>
</html>
