//=========变量定义区域===========================================================


//==================动作响应区域==================================================
//页面初始化
function pageInit(){
	
}

function onNewPsw(){
	modaltext = '<form id="useredit" action="" class="form-inline" role="form" style="margin:10px 0px 10px 60px" method="post">';
	
	modaltext +='<div class="row col-sm-12 form-group"><label for="OldPswd">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;旧密码:</label>';
	modaltext +='<input type="password" class="form-control" id="OldPswd" name="OldPswd" placeholder="请输入旧密码" required onchange="onOldPswdChange()">';
	modaltext +='</div>';
	
	modaltext +='<div class="row col-sm-12 form-group"><label for="NewPswd" >&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;新密码:</label>';
	modaltext +='<input type="password" class="form-control" id="NewPswd" name="NewPswd" placeholder="请输入新密码" required onblur="onPswd();">';
	modaltext +='</div>';
	
	modaltext +='<div class="row col-sm-12 form-group"><label for="NewPswd2" class="control-label">重复新密码:</label>';
	modaltext +='<input type="password" class="form-control" id="NewPswd2" name="NewPswd2" placeholder="请再输入一遍新密码" required onblur="onPswd();">';
	modaltext +='<span class="text-danger" id="msg" style="display:none"></span></div>';
	
	modaltext +='<div class="row col-sm-12 ">';
	modaltext +='<div class="mx-auto"><button name="submit" class="btn btn-success form-control" id="submit" disabled onclick="onSubmit()">提交</button><span id="backmsg"></span>';
	modaltext +='</div></form>';
	ShowModal("修改密码",modaltext);
}
function onSubmit(){
	var old = $("#OldPswd").val();
	var new1 = $("#NewPswd").val();
	var new2 = $("#NewPswd2").val();
	$.post("api/updatepswd",{OldPswd:old,NewPswd:new1,NewPswd2:new2},function(data){
        var msg='';
		if (data==1){
			msg='成功:密码修改成功!';
		}else{
			msg='失败:密码修改失败:'+data;
		};
		alert(msg);
    });
	
}
function onPswd(){//新密码改变
	var p1 = $("#NewPswd").val();
	var p2 = $("#NewPswd2").val();
	if(p1!=p2){//两次密码不同
		$("#msg").text("两次输入的密码不同！");
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
	ShowModal("提示",md5('你好'));
}
$(document).ready(function (){
	$("#exit").before('退出')
	$("#editPd").html("修改密码")
	$('#editMs').html("编辑信息")
	$('#userInformation').html("用户基本信息表")
	$("#loginName").html('登 录 名')
	$("#name").html('姓    名')
	$("#IDcard").html('身份证号')
	$("#phoneNumber").html("手 机 号")
	$("#email").html("邮    箱")
	$("#sex").html("性    别")
	$("#qq").html("QQ 号 码")
	$("#wechat").html("微信号码")
});

//=========AJAX函数定义区域=======================================================
