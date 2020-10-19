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


$(document).ready(function(){
	$("#webApi").html('WebAPI接口文档');
	$("#granafa").html('Grafana读取平台数据接口文档');
	$("#exit").after('退出');
});