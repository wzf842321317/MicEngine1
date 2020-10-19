function zTreeOnClick(event, treeId, treeNode) {
	$("#SelectNodeName").text(treeNode.name);
	$("#SelectNodeLevelCode").text(treeNode.treelevelcode);
};
//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//页面初始化
function pageInit(){

}
//获取界面的cookie 进行判断切换界面
//获取cookie

function setCookie(c_name, value) {
	var exdate = new Date();
	document.cookie = c_name + "=" + escape(value);
}

$(document).ready(function(){
	$("#webApi").html("WebAPI Interface documentation");
	$("#granafa").html("Grafana Read platform data interface document");
	$("#exit").after('Exit');
});
