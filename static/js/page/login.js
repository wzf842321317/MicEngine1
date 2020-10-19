function onCaptchaInput(){
   	var input_chkcod=$('#chkcod').val();
    var captcha=$('#captcha').val();
   	if(input_chkcod == captcha){
        $('#submit').removeAttr('disabled');
		$('#submit').attr('class','w-100 btn btn-success')
    }else{
    	$('#submit').attr('disabled',"true");
		$('#submit').attr('class','w-100 btn btn-outline-success')
    }
}

function onPswdChange(){
	var pwd=$("#password").val();
	var pwd=$("#password").val(MD5(pwd));
}

$(document).ready(function(){
	$("#tip").val(window.location.host);
	$("#name").attr('placeholder','请输入帐号');
	$("#password").attr('placeholder','请输入登录密码');
	$("#chkcod").attr('placeholder','请输验证码');
	$("#submit").text('登录');

});
