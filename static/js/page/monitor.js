//=========变量定义区域===========================================================


//==================动作响应区域==================================================
function zTreeOnClick(event, treeId, treeNode) {
	if (treeNode.href.length > 0){
		var width = $("#MonitorFrame").innerWidth();
		var winH = $(window).height();
		$("#MonitorView").attr("width",width);
		$("#MonitorView").attr("height",(winH-80));
		$("#MonitorView").attr("src",treeNode.href);
	}
}
//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//页面初始化
function pageInit(){
	var width = $("#MonitorFrame").innerWidth();
	var winH = $(window).height();
	$("#MonitorView").attr("width","105%");
	$("#MonitorView").attr("height",(winH-10));
}
$(document).ready(function() {
	//Ztree国际化
	$("#exit").after('退出');
	$("#ExpandTreeNode").html('展开');
	$("#CollapseTreeNode").html('折叠');
	$("#HideTreeNode").html('隐藏');
	$("#SearchTreeNode").attr("placeholder", '搜索');

	$('#ExpandTreeNode').attr("title", "未选中节点时展开所有节点,选中节点时展开选中节点");
	$('#CollapseTreeNode').attr("title", "未选中节点时折叠所有节点,选中节点时折叠选中节点");
});


function onMonitorTreeNodeView(isview){
	if (isview==1){
	  $("#TreeBar").show();
	  $("#TreeBarSm").hide();
	  $("#FloatArea").attr("class","col-sm-9 col-md-6 col-lg-8 col-xl-10");
	}else{
	  $("#TreeBar").hide();
	  $("#TreeBarSm").show();
	  $("#FloatArea").attr("class","col-12");
	}
	$("#MonitorView").attr("width","105%");
}
//=========AJAX函数定义区域=======================================================