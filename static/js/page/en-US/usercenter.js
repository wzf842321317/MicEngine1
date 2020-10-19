//=========变量定义区域===========================================================


//==================动作响应区域==================================================
//页面初始化
function pageInit(){
	
}

function onNewPsw(){
	modaltext = '<form id="useredit" action="" class="form-inline" role="form" style="margin:10px 0px 10px 60px" method="post">';
	
	modaltext +='<div class="row col-sm-12 form-group"><label for="OldPswd">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Old password:</label>';
	modaltext +='<input type="password" class="form-control" id="OldPswd" name="OldPswd" placeholder="Please enter the old password" required onchange="onOldPswdChange()">';
	modaltext +='</div>';
	
	modaltext +='<div class="row col-sm-12 form-group"><label for="NewPswd" >&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;New password:</label>';
	modaltext +='<input type="password" class="form-control" id="NewPswd" name="NewPswd" placeholder="Please enter a new password" required onblur="onPswd();">';
	modaltext +='</div>';
	
	modaltext +='<div class="row col-sm-12 form-group"><label for="NewPswd2" class="control-label">Repeat new password:</label>';
	modaltext +='<input type="password" class="form-control" id="NewPswd2" name="NewPswd2" placeholder="Please enter the new password again" required onblur="onPswd();">';
	modaltext +='<span class="text-danger" id="msg" style="display:none"></span></div>';
	
	modaltext +='<div class="row col-sm-12 ">';
	modaltext +='<div class="mx-auto"><button name="submit" class="btn btn-success form-control" id="submit" disabled onclick="onSubmit()">Submit</button><span id="backmsg"></span>';
	modaltext +='</div></form>';
	ShowModal("Change Password",modaltext);
}
function onSubmit(){
	var old = $("#OldPswd").val();
	var new1 = $("#NewPswd").val();
	var new2 = $("#NewPswd2").val();
	$.post("api/updatepswd",{OldPswd:old,NewPswd:new1,NewPswd2:new2},function(data){
        var msg='';
		if (data==1){
			msg='Success: password changed successfully!';
		}else{
			msg='Failed: password modification failed:'+data;
		};
		alert(msg);
    });
	
}
function onPswd(){//新密码改变
	var p1 = $("#NewPswd").val();
	var p2 = $("#NewPswd2").val();
	if(p1!=p2){//两次密码不同
		$("#msg").text("The two passwords are different！");
		$("#msg").show();//显示提示信息
		$("#submit").attr("disabled",true);//禁用提交按钮
	}else{//两次密码相同
		$("#submit").attr("disabled",false);//启用提交按钮
		$("#NewPswd").val(MD5(p1))//加密
		$("#NewPswd2").val(MD5(p2))//加密
		$("#msg").hide();//隐藏提示信息
	}
}
//旧密码改变,加密
function onOldPswdChange(){
	var pswd = $("#OldPswd").val();
	$("#OldPswd").val(MD5(pswd));
}

function onEditUserMsg(){
	ShowModal("Tips",md5('Hello'));
}

$(document).ready(function (){
	$("#exit").before('Exit')
	$("#editPd").html("Change Password")
	$('#editMs').html("Edit information")
	$('#userInformation').html("User basic information table")
	$("#loginName").html('Login name')
	$("#name").html('Name')
	$("#IDcard").html('ID number')
	$("#phoneNumber").html("Phone number")
	$("#email").html("Email")
	$("#sex").html("Gender")
	$("#qq").html("QQ number")
	$("#wechat").html("Wechat number")
});

//=========AJAX函数定义区域=======================================================
