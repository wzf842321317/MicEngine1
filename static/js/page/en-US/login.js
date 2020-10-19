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
	$("#name").attr('placeholder','Please Input Your Username');
	$("#password").attr('placeholder','Please Input Your Password');
	$("#chkcod").attr('placeholder','Please Input Verification Code');
	$("#submit").text('Sign in for MicEngine');
});
