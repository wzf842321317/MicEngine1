<SCRIPT type="text/javascript">
//树的配置
var zTreeSetting = {
	//事件回调函数
	callback:{
		onClick:zTreeOnClick,//用于捕获节点被点击的事件回调函数
    onDblClick:zTreeOnDbClick//用于捕获节点被双击击的事件回调函数
	},
	//数据格式的设置
	data:{
		simpleData:{
			enable:true,  //使用简单的json数据
			idKey:"id",   //节点数据中保存唯一标识的属性名称
			pIdKey:"pId", //节点数据中保存其父节点唯一标识的属性名称
			rootPId:{{.}}  //用于修正根节点父节点数据，即 pIdKey 指定的属性值
		}
	},
	//树的显示设置
	view:{
		showLine:true, //设置是否显示连线 默认为true
		showTitle: true,//设置是否显示鼠标悬浮时显示title的效果,请务必与 setting.data.key.title 同时使用。
		showIcon:true,//是否显示节点的图标
		dblClickExpand: true,//设置是否支持双击展开树节点
		fontCss : {color:"blue"},//设置节点样式,例如：{color:"#ff0011", background:"blue"}  
		expandSpeed: "fast",//设置展开的速度  fast normal  slow 
		nameIsHTML: true//名字是否是HTML
	}
};

function expandNode(e) {
	var zTree = $.fn.zTree.getZTreeObj("NodeTree"),
	type = e.data.type,
	nodes = zTree.getSelectedNodes();

  if (type=="expand"){
    if (nodes.length==0)//如果没有选中节点
      type = "expandAll";//展开所有
    else//有选中
      type = "expandSon";//展开选中(含子节点)
  }
  if (type=="collapse"){
    if (nodes.length==0)//如果没有选中节点
      type = "collapseAll";//折叠所有
    else//有选中
      type = "collapseSon";//折叠选中(含子节点)
  }
  
	if (type == "expandAll") {//展开全部节点
		zTree.expandAll(true);
	} else if (type == "collapseAll") {//折叠全部节点
		zTree.expandAll(false);
	} else {
		var callbackFlag = $("#callbackTrigger").attr("checked");
		for (var i=0, l=nodes.length; i<l; i++) {
			zTree.setting.view.fontCss = {};
			if (type == "expand") {//展开单个节点
				zTree.expandNode(nodes[i], true, null, null, callbackFlag);
			} else if (type == "collapse") {//折叠单个节点
				zTree.expandNode(nodes[i], false, null, null, callbackFlag);
			} else if (type == "toggle") {//展开折叠切换单个节点
				zTree.expandNode(nodes[i], null, null, null, callbackFlag);
			} else if (type == "expandSon") {//展开单个节点(含子节点)
				zTree.expandNode(nodes[i], true, true, null, callbackFlag);
			} else if (type == "collapseSon") {//折叠单个节点(含子节点)
				zTree.expandNode(nodes[i], false, true, null, callbackFlag);
			}
		}
	}
}
var winH = $(window).height();
$(document).ready(function(){
	//初始化zTree
	$.fn.zTree.init($("#NodeTree"), zTreeSetting, zNodes);
  fuzzySearch('NodeTree','#SearchTreeNode',false,false); //初始化模糊搜索方法
	$("#ExpandTreeNode").bind("click", {type:"expand"}, expandNode);
	$("#CollapseTreeNode").bind("click", {type:"collapse"}, expandNode);
  
  //初始化树区域的最大高度
  if(winH < 600){
    winH=600;
  } 
  $("#NodeTree").attr("style","height:"+(winH-70)+"px; max-height:"+(winH-70)+"px;overflow: auto;");
  $("#FloatArea").attr("style","height:"+(winH-70)+"px; max-height:"+(winH)+"px;overflow: auto;");
  
  //其他初始化工作
  pageInit();
});

function onTreeNodeView(isview){
  if (isview==1){
    $("#TreeBar").show();
    $("#TreeBarSm").hide();
    $("#FloatArea").attr("class","col-sm-9 col-md-6 col-lg-8 col-xl-10");
  }else{
    $("#TreeBar").hide();
    $("#TreeBarSm").show();
    $("#FloatArea").attr("class","col-12");
  }
}
</SCRIPT>	