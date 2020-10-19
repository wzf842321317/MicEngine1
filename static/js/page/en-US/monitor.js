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
	$("#exit").after('Exit');
	$("#ExpandTreeNode").html('Open');
	$("#CollapseTreeNode").html('Fold');
	$("#HideTreeNode").html('Hide');
	$("#SearchTreeNode").attr("placeholder", 'Search');
	$('#ExpandTreeNode').attr("title", "Expand all nodes when no nodes are selected, and expand selected nodes when nodes are selected");
	$('#CollapseTreeNode').attr("title", "Collapse all nodes when no nodes are selected, and collapse selected nodes when nodes are selected");

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