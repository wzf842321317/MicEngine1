<!DOCTYPE html>
<html lang="cn">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<title>{{.WebTitle}}</title>
<link rel="Bookmark" href="{{.IcoPic}}" /> 
<link rel="shortcut icon" href="{{.IcoPic}}" /> 
<link rel="stylesheet" href="static/bootstrap-4.4.1/css/bootstrap.css" type="text/css">
<link href="static/css/login.css" type="text/css" rel="stylesheet">

</head>
<body>
	<div class="login">
	  <div class="message">{{.ProjectName}}</div>
	  <div id="darkbannerwrap"></div>
    <form name="login" action="login" method="post">
			<input  title="" name="Name" id="name"  placeholder="" autocomplete="off" type="text" required maxlength="16">
				<hr class="hr15">
			<input title="" name="Password" id="password" placeholder="" autocomplete="off" type="password" required maxlength="16" onchange="onPswdChange()">
				<hr class="hr15">
			<input title="" name="chkcod" id="chkcod" placeholder="" type="text" required maxlength="16" oninput="onCaptchaInput()">
			<img src=static/img/code/{{.Chkcod_1}}.gif>
      <img src=static/img/code/{{.Chkcod_2}}.gif>
      <img src=static/img/code/{{.Chkcod_3}}.gif>
      <img src=static/img/code/{{.Chkcod_4}}.gif>
      {{.Sess}}
      <hr class="hr15">
			<input name="captcha" id="captcha" type="hidden" value={{.Captcha}} >
			<input title="浏览器IP" name="tip" id="tip" type="text" hidden>
             <button type="submit" class="w-100 btn btn-outline-success" name="submit" id="submit" disabled="true"></button>
			  <hr class="hr20">
		</form>
	</div>
  <script src='static/js/jquery-3.4.1.min.js'></script>
  <script type="text/javascript" src="static/js/PageTpl.js"></script>  
  <script type="text/javascript" src="static/bootstrap-4.4.1/js/bootstrap.min.js"></script> 
  <script src='static/js/page/{{.JsFileName}}login.js'></script>
</body>
</html>