var LAST_SELECTED_CELL='A1';//上次选择的单元格
var CELLSMAP=new Array();//单元格Map,以单元格坐标字符串为key

//单元格被单击
function onClickCell(axis){
    var msg = CELLSMAP[axis].Value;//获取单元格的内容
    setCellStatus(axis,true);
    if(LAST_SELECTED_CELL!=axis){//如果与上次单击的不是同一个单元格
        setCellStatus(LAST_SELECTED_CELL,false);
    }else{//与上次单击的是同一个单元格
        setCellStatus(LAST_SELECTED_CELL,true);      
    }
    if(isMicFormFormula(msg)==true){
        $("#CheckCellValue").attr("class","btn btn-success btn-sm");//校验按钮样式
        $("#CheckCellValue").removeAttr("disabled");//校验按钮使能
    }else{
        $("#CheckCellValue").attr("class","btn btn-outline-success btn-sm");//校验按钮样式
        $("#CheckCellValue").attr("disabled","disabled"); //校验按钮禁用
    }
    
    $("#selectTdValue").val(msg);//显示单元格的内容
	$("#coordinate").html(axis);//显示单元格坐标
    $("#CellAxis").val(axis);//显示单元格坐标
    if(CELLSMAP[axis].Formula.length>0){//如果选中的单元格中有Excel公式
		$("#selectTdValue").val("="+CELLSMAP[axis].Formula);//显示Excel公式
    }
    setTimeout(function (){//延时0.5秒变更上次选择的单元格坐标
        LAST_SELECTED_CELL = axis;//记录选中的单元格坐标
    }, 200);
    
}

//单元格的值发生改变
function onCellValueChange(){
    var cellinput = $("#selectTdValue").val();//获取输入的值
    CELLSMAP[LAST_SELECTED_CELL].Status=1;//记录状态
    if (cellinput.length>2){//如果输入的值长度大于2
        if(cellinput[0]=="="){//判断是否是Excel公式
            if(CELLSMAP[LAST_SELECTED_CELL].Value.length>0){//如果有值
                onPostSetCellValue(LAST_SELECTED_CELL,"");//清除值
            }
            $("#"+LAST_SELECTED_CELL).text("");//刷新显示
            var formula=cellinput.slice(1);//如果是公式,获取公式的内容(=号之后的内容)
            CELLSMAP[LAST_SELECTED_CELL].Formula = formula;//保存公式
            onPostSetCellFormula(LAST_SELECTED_CELL,formula);//保存单元格的公式
        }else{//如果不是Excel公式
            if(CELLSMAP[LAST_SELECTED_CELL].Formula.length>0){//如果有Excel公式
                onPostSetCellFormula(LAST_SELECTED_CELL,"");//清除公式
            }
            if (isMicFormFormula(cellinput)==true){//包含脚本计算公式
                $("#CheckCellValue").attr("class","btn btn-success btn-sm");//校验按钮样式
                $("#CheckCellValue").removeAttr("disabled");//校验按钮使能
                var formulas = getMicFormFormula(cellinput);//获取脚本
                compileMicFormFormula(formulas,cellinput);
            }else{//不包含脚本计算公式
                $("#CheckCellValue").attr("class","btn btn-outline-success btn-sm");//校验按钮样式
                $("#CheckCellValue").attr("disabled","disabled"); //校验按钮禁用
                $("#"+LAST_SELECTED_CELL).text(cellinput);
                onPostSetCellValue(LAST_SELECTED_CELL,cellinput);//保存单元格的值
            }
            $("#"+LAST_SELECTED_CELL).text(cellinput);//刷新显示
        }
    }else{//如果输入的值长度不大于2
        if(CELLSMAP[LAST_SELECTED_CELL].Formula.length>0){//如果有Excel公式
            onPostSetCellFormula(LAST_SELECTED_CELL,"");//清除公式
        }
        onPostSetCellValue(LAST_SELECTED_CELL,cellinput);//保存单元格的值
        $("#"+LAST_SELECTED_CELL).text(cellinput);//刷新显示
    }
    setCellStatus(LAST_SELECTED_CELL,true);//设置显示状态
}
//单元格的值有输入时
function onCellValueInput(){
    var cellinput = $("#selectTdValue").val();//获取输入的值
    if(isMicFormFormula(cellinput)==true){
        $("#CheckCellValue").attr("class","btn btn-success btn-sm");//校验按钮样式
        $("#CheckCellValue").removeAttr("disabled");//校验按钮使能
    }else{
        $("#CheckCellValue").attr("class","btn btn-outline-success btn-sm");//校验按钮样式
        $("#CheckCellValue").attr("disabled","disabled"); //校验按钮禁用
    }
}

//单元格校验按钮
function onCheckCellValue(){
    var cellinput = $("#selectTdValue").val();//获取输入的值
    var formulas = getMicFormFormula(cellinput);//获取脚本
    var formData = new FormData();
    var formula="";
    for(var i=0;i<formulas.length;i++){
        formula +=formulas[i];
        if(i<(formulas.length-1)){
            formula+=";";
        } 
    }
    formData.append("script",formula);
    
	var xhr=new XMLHttpRequest();
    xhr.open("post","api/script/compile");
    xhr.send(formData);
    xhr.onload=function(){
        if(xhr.status==200){
			if(xhr.responseText!="ok"){
                CELLSMAP[LAST_SELECTED_CELL].Status=2;
                setCellStatus(LAST_SELECTED_CELL,true);
                ShowModal("Compilation error",'<div class="alert alert-danger"><strong>Error:</strong><br>'+xhr.responseText+'</div>');
            }else{
                CELLSMAP[LAST_SELECTED_CELL].Status=1;
                setCellStatus(LAST_SELECTED_CELL,true);
                ShowModal("Compilation passed",'<div class="alert alert-success"><strong>Success:</strong>Compilation passed！</div>');
                onPostSetCellValue(LAST_SELECTED_CELL,cellinput);//保存单元格的值
            }
        }
    }
}

//校验编译报表公式
function compileMicFormFormula(formulas,cellinput){
    var formData = new FormData();
    var formula="";
    for(var i=0;i<formulas.length;i++){
        formula +=formulas[i];
        if(i<(formulas.length-1)){
            formula+=";";
        } 
    }
    console.log(formula);
    formData.append("script",formula);
    
	var xhr=new XMLHttpRequest();
    xhr.open("post","api/script/compile");
    xhr.send(formData);
    xhr.onload=function(){
        if(xhr.status==200){
			if(xhr.responseText!="ok"){
                CELLSMAP[LAST_SELECTED_CELL].Status=2;
                setCellStatus(LAST_SELECTED_CELL,true);
                ShowModal("Compilation error",'<div class="alert alert-danger"><strong>Error:</strong><br>'+xhr.responseText+'</div>');
            }else{
                CELLSMAP[LAST_SELECTED_CELL].Status=1;
                setCellStatus(LAST_SELECTED_CELL,true);
                onPostSetCellValue(LAST_SELECTED_CELL,cellinput);//保存单元格的值
            }
        }
    }
}

//是否包含报表公式
function isMicFormFormula(cellmsg){
    var patt=/\{\{[^\{\}]+\}\}/g;//新建正则表达式
    return patt.test(cellmsg);
}

//提取报表公式
function getMicFormFormula(cellmsg){
    var patt=/\{\{[^\{\}]+\}\}/g;//新建正则表达式
    formulas=cellmsg.match(patt);
    for(var i=0;i<formulas.length;i++){
        var formula = formulas[i];
        formula = formula.substr(2,formula.length-4)
        formulas[i] = formula;
    }
    console.log(formulas);
    return formulas;
}

//设置单元格的状态样式
function setCellStatus(axis,withborder){
    if(withborder==false){
        switch (CELLSMAP[axis].Status){
            case 0://单元格内容没有改变
                $("#"+axis).removeAttr("style");
                break;
            case 1://单元格内容有改变
                $("#"+axis).attr("style","background-color:#99ff00");
                break;
            case 2://单元格公式有改变,但没有通过校验
                $("#"+axis).attr("style","background-color:#ff0066");
                break;
            default://单元格公式有改变,但没有保存成功
                $("#"+axis).attr("style","background-color:#ffff00");
                break;
        }
    }else{
        switch (CELLSMAP[axis].Status){
            case 0://单元格内容没有改变
                $("#"+axis).attr("style","border:blue solid 2px;");
                break;
            case 1://单元格内容有改变
                $("#"+axis).attr("style","border:blue solid 2px;background-color:#99ff00");
                break;
            case 2://单元格公式有改变,但没有通过校验
                $("#"+axis).attr("style","border:blue solid 2px;background-color:#ff0066");
                break;
            default://单元格公式有改变,但没有保存成功
                $("#"+axis).attr("style","border:blue solid 2px;background-color:#ffff00");
                break;
        } 
    }
}

//保存单元格的值
function onPostSetCellValue(axis,value){
    var formData = new FormData();

	formData.append("filepath",NODES_KVA[REPORT_ID].TemplateUrl);
	formData.append("filename",NODES_KVA[REPORT_ID].TemplateFile);
	formData.append("cellaxis",axis);
    formData.append("cellvalue",value);

	var xhr=new XMLHttpRequest();
    xhr.open("post","api/setcellvalue");
    xhr.send(formData);
    xhr.onload=function(){
        if(xhr.status==200){
            if(xhr.responseText!="ok"){
                CELLSMAP[axis].Status = 3;
                setCellStatus(axis,axis==LAST_SELECTED_CELL?true:false)
                ShowModal("Save file error",'<div class="alert alert-danger"><strong>Error:</strong><br>'+xhr.responseText+'</div>');
            }else{
                CELLSMAP[axis].Value=value;
                $("#"+axis).text(value);//刷新显示
            }
        }
    }
}

//保存单元格的公式
function onPostSetCellFormula(axis,formula){
    var formData = new FormData();

	formData.append("filepath",NODES_KVA[REPORT_ID].TemplateUrl);
	formData.append("filename",NODES_KVA[REPORT_ID].TemplateFile);
	formData.append("cellaxis",axis);
    formData.append("cellformula",formula);

	var xhr=new XMLHttpRequest();
    xhr.open("post","api/setcellformula");
    xhr.send(formData);
    xhr.onload=function(){
        if(xhr.status==200){
			if(xhr.responseText!="ok"){
                CELLSMAP[axis].Status = 3;
                setCellStatus(axis,axis==LAST_SELECTED_CELL?true:false)
                ShowModal("Save file error",'<div class="alert alert-danger"><strong>Error:</strong><br>'+xhr.responseText+'</div>');
            }else{
                CELLSMAP[axis].Formula=formula;
            }
        }
    }
}
