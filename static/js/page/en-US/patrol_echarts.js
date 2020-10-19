//<!--Echarts脚本文件-->
var Echarts_dom = document.getElementById("Modal_Echarts");
var DataEcharts = echarts.init(Echarts_dom);

var pageEcharts = echarts.init(document.getElementById("PageEcharts"));
//------------饼图----------------------
echarts_option = {
    color: ['#ff7f50','#87cefa','#da70d6','#32cd32','#6495ed',
    '#ff69b4','#ba55d3','#cd5c5c','#ffa500','#40e0d0',
    '#1e90ff','#ff6347','#7b68ee','#00fa9a','#ffd700',
    '#6699FF','#ff6666','#3cb371','#b8860b','#30e0e0'],
};


function refreshEcharts(group)
{	
   DataEcharts.setOption(echarts_option, true);
    DataEcharts.setOption({
        tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b} : {c} ({d}%)'
        },
		title: {
            text: 'Inspection workload statistics',
            subtext: 'Statistics by person',
            left: 'center'
        },
        series:[
            {
                name: 'Inspector',
                type: 'pie',
                radius: '70%',
                center: ['50%', '60%'],
                data: function(){
                    var serie=[];
                    for (let key in group) {
                           var item={
                            value:group[key],
                            name:key
                        };
                        serie.push(item);
                    };
                    return serie;
                }(),
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    });
    
	$("#Modal_Echarts").show();
};

function refreshTrendEcharts(tagname,checkdatas)
{	
   DataEcharts.setOption(echarts_option, true);
   var times=[];
   var datas=[];
   //console.log(checkdatas);
   var i = checkdatas.length, j;
   var tempExchangVal;
   while (i > 0) {
       for (j = 0; j < i - 1; j++) {
           var t1=new Date(checkdatas[j].time.replace(/T/," "));//基准时间;
           var t2=new Date(checkdatas[j+1].time.replace(/T/," "));//基准时间;
           if (t1 > t2) {
               tempExchangVal = checkdatas[j];
               checkdatas[j] = checkdatas[j + 1];
               checkdatas[j + 1] = tempExchangVal;
           }
       }
       i--;
   }
   for(var i=0;i<checkdatas.length;i++){
       times.push(checkdatas[i].time);
       datas.push(checkdatas[i].value);
   }
    DataEcharts.setOption({
        tooltip: {
            trigger: 'axis',
            axisPointer: {            // 坐标轴指示器，坐标轴触发有效
                type: 'line'        // 默认为直线，可选为：'line' | 'shadow'
            }
        },
        toolbox: {
            show : true,
            feature : {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                restore: {},
                dataView : {show: true, readOnly: false},
                magicType : {show: true, type: ['line', 'bar']},
                restore : {show: true},
                saveAsImage : {show: true}
            }
        },
		title: {
            text: tagname,
            subtext: '',
            left: 'center'
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        xAxis: [
            {
                type: 'category',
                data: times,
                axisTick: {
                    alignWithLabel: true
                }
            }
        ],
        yAxis: [
            {
                type: 'value'
            }
        ],
        series: [
            {
                markPoint : {
                    data : [
                        {type : 'max', name : '最大值'},
                        {type : 'min', name : '最小值'}
                    ]
                },
                markLine : {
                    data : [
                        {type : 'average', name : '平均值'}
                    ]
                },
                step: 'end',
                name: tagname,
                type: 'line',
                data: datas
            }
        ]
    });
};

function refreshRelitionEcharts(datas){
    DataEcharts.setOption(echarts_option, true);
    DataEcharts.setOption({
        tooltip: {
            trigger: 'item',
            triggerOn: 'mousemove'
        },
        series:[
            {
                type: 'tree',
   
                data: [datas],
    
                top: '1%',
                left: '15%',
                bottom: '1%',
                right: '20%',
    
                symbolSize: 7,
    
                label: {
                    position: 'left',
                    verticalAlign: 'middle',
                    align: 'right'
                },
    
                leaves: {
                    label: {
                        position: 'right',
                        verticalAlign: 'middle',
                        align: 'left'
                    }
                },
    
                expandAndCollapse: true,
    
                animationDuration: 550,
                animationDurationUpdate: 750
            }
        ]
    });    
}

function refreshPageEcharts(sites,datas){
    pageEcharts.setOption(echarts_option, true);
    pageEcharts.setOption({
        tooltip: {
            trigger: 'item',
            triggerOn: 'mousemove'
        },
        title: {
            text: "Inspection completion rate",
            subtext: 'Statistics by checkpoint',
            left: 'center'
        },
        toolbox: {
            show : true,
            feature : {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                restore: {},
                dataView : {show: true, readOnly: false},
                magicType : {show: true, type: ['line', 'bar']},
                restore : {show: true},
                saveAsImage : {show: true}
            }
        },
        xAxis: {
            type: 'category',
            data: sites
        },
        yAxis: {
            type: 'value'
        },
        series: [{
            data: datas,
            type: 'bar'
        }]
    });    
    $("#PageEcharts").show();
}

if (echarts_option && typeof echarts_option === "object") {
    DataEcharts.setOption(echarts_option, true);
	document.getElementById("Modal_Echarts").style.display="none";//隐藏
};
if (echarts_option && typeof echarts_option === "object") {
    DataEcharts.setOption(echarts_option, true);
	document.getElementById("PageEcharts").style.display="none";//隐藏
};