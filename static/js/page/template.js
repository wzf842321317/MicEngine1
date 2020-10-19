//=========变量定义区域===========================================================
var SAMPLE_NODE={//当前选中的样本节点
	Id:0,//样本ID
	Pid:0,//所在作业ID
	Name:"",//样本名称
	PoolId:0,//所属样本池ID
	TreeLevel:"",//层级码
	SamplingSite:"",//取样地点
	FuncType:0,//样本类型码
	FuncName:"",//样本类型名
	IsRegular:0,//是否常规样/临时样
	BaseTime:"",//所在车间基准时间
	ShiftHour:0//所在车间每班工作时间
};
var HAS_SELECT_SAMPLE=false;//已经选择了样本节点
var SAMPLE_LAB_TAGS=[];//化验标签信息
var SAMPLE_LAB_RESULT=[];//化验结果信息

//==================动作响应区域==================================================
//响应鼠标单击
function zTreeOnClick(event, treeId, treeNode) {
	if (treeNode.nodetype == 9999){

	}
}

//响应鼠标双击
function zTreeOnDbClick(event, treeId, treeNode){
	
}
//页面初始化工作
function pageInit(){

};
//=========AJAX请求定义区域=======================================================

//=========AJAX加载定义区域=======================================================

//=========AJAX数据接收解析区域====================================================