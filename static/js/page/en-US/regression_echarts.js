var his_dom = document.getElementById("Echarts_His");
var hisChart = echarts.init(his_dom); 

var hisInterval_dom = document.getElementById("Echarts_HisIntervalSerial");
var hisIntervalChart = echarts.init(hisInterval_dom); 

var scatter_dom = document.getElementById("Echarts_Scatter");
var scatterChart = echarts.init(scatter_dom); 

var trend_dom = document.getElementById("Echarts_Trend");
var trendChart = echarts.init(trend_dom); 
//------------单变量历史趋势图(通用型)----------------------
 his_option = {
    title: {
        text: '原始历史数据'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross'
        }
    },

	color:['red','green','blue','#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83',  '#ca8622', '#bda29a','#6e7074', '#546570', '#c4ccd3'],
    legend: {
        data: ['Raw data','Increment']
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
	grid: {
        left: '5%',
        right: '6%'
		//containLabel: true
    },
    xAxis: [{
        type: 'category',
		name: 'Time',
        boundaryGap: false
        //data:  [1,2,3,4,5,6,7,8]
    }],
    yAxis: [{
		gridIndex: 0,
		name: 'Numerical value',
        type: 'value',
		position: 'left'
    }],
	dataZoom: [{//缩放条
        type: 'inside',
        start: 0,
        end: 100
    }, {
        start: 0,
        end: 100,
        handleIcon: 'M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z',
        handleSize: '80%',
        handleStyle: {
            color: '#fff',
            shadowBlur: 3,
            shadowColor: 'rgba(0, 0, 0, 0.6)',
            shadowOffsetX: 2,
            shadowOffsetY: 2
        }
    }],
	series: [
        {
            name:'Raw data',
            type:'line',
			markPoint : {
                data : [
                    {type : 'max', name : 'Max'},
                    {type : 'min', name : 'Min'}
                ]
            },
            markLine : {
                data : [
                    {type : 'average', name : 'Average'}
                ]
            },
			step: '',
            data: [1,2,3,4,5,6,7,8]
        },{
            name:'Increment',
            type:'line',
			markPoint : {
                data : [
                    {type : 'max', name : 'Max'},
                    {type : 'min', name : 'Min'}
                ]
            },
            markLine : {
                data : [
                    {type : 'average', name : 'Average'}
                ]
            },
			step: '',
            data: [1,2,3,4,5,6,7,8]
        }
    ]
}; 

//------------单变量历史趋势图(时间对齐)----------------------
 hisInterval_option = {
    title: {
        text: 'Aggregate historical data at equal intervals'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross'
        }
    },

	color:['blue','red','green','#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83',  '#ca8622', '#bda29a','#6e7074', '#546570', '#c4ccd3'],
    legend: {
        data: ['Aggregate data at equal intervals']
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
	grid: {
        left: '5%',
        right: '6%'
		//containLabel: true
    },
    xAxis: [{
        type: 'category',
		name: 'Time',
        boundaryGap: false,
        data:  [1,2,3,4,5,6,7,8]
    }],
    yAxis: [{
		gridIndex: 0,
		name: 'Numerical value',
        type: 'value',
		position: 'left'
    }],
	dataZoom: [{//缩放条
        type: 'inside',
        start: 0,
        end: 100
    }, {
        start: 0,
        end: 100,
        handleIcon: 'M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z',
        handleSize: '80%',
        handleStyle: {
            color: '#fff',
            shadowBlur: 3,
            shadowColor: 'rgba(0, 0, 0, 0.6)',
            shadowOffsetX: 2,
            shadowOffsetY: 2
        }
    }],
	series: [
        {
            name:'Aggregate data at equal intervals',
            type:'line',
			markPoint : {
                data : [
                    {type : 'max', name : 'Max'},
                    {type : 'min', name : 'Min'}
                ]
            },
            markLine : {
                data : [
                    {type : 'average', name : 'Average'}
                ]
            },
			step: '',
            data: [1,2,3,4,5,6,7,8]
        }
    ]
}; 
//------------散点图----------------------
 scatter_option = {
	title: {
        text: 'Scatter plot of linear regression data',
        left: 'center'
    },
	color:['red','green','blue','#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83',  '#ca8622', '#bda29a','#6e7074', '#546570', '#c4ccd3'],
	toolbox: {
        show : true,
        feature : {
			dataZoom: {
                yAxisIndex: 'none'
            },
            restore: {},
            dataView : {show: true, readOnly: false},
            restore : {show: true},
            saveAsImage : {show: true}
        }
    },
	xAxis: {
		name: 'Actual value',
        scale: true
    },
    yAxis: {
		name: 'Estimate',
        scale: true
    },
	tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross'
        }
    },
    series: [{
		type : 'line',//中心线
		lineStyle: {
            normal: {
                type: 'dashed'
            }
        },
		data: [[0,0],[100,100]]
	},{
		type : 'line',//上西格玛线
		lineStyle: {
            normal: {
                type: 'dashed'
            }
        },
		data: [[0,0],[100,100]]
	},{
		type : 'line',//下西格玛线
		lineStyle: {
            normal: {
                type: 'dashed'
            }
        },
		data: [[0,0],[100,100]]
	},{
        type: 'scatter',
		symbolSize: 20,
        data: [[161.2, 51.6], [167.5, 59.0], [159.5, 49.2], [157.0, 63.0], [155.8, 53.6],
            [170.0, 59.0], [159.1, 47.6], [166.0, 69.8], [176.2, 66.8], [160.2, 75.2],
            [172.5, 55.2], [170.9, 54.2], [172.9, 62.5], [153.4, 42.0], [160.0, 50.0]
        ],
    }]
};
//------------折线趋势图----------------------
trend_option = {
    title: {
        text: 'Data trend chart'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross'
        }
    },
	//贯穿的时间纵线
	axisPointer: {
        link: {xAxisIndex: 'all'}
    },
	color:['red','green','blue','#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83',  '#ca8622', '#bda29a','#6e7074', '#546570', '#c4ccd3'],
    legend: {
        data:['Actual value', 'Model value', 'Residual value', 'Standard residual', 'Relative deviation']
    },
	//网格线位置控制
    grid: [{
        left: 40,
        right: 40,
        height: '40%',
		//containLabel: true
    }, {
        left: 40,
        right: 40,
        top: '55%',
        height: '40%',
		//containLabel: true
    }],
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
    xAxis: [{
        type: 'category',
        boundaryGap: false,
        data: [0,1,2,3]
    },{
		gridIndex: 1,
        type: 'category',
        boundaryGap: false,
        data: [0,1,2,3]
    }],
    yAxis: [{
		gridIndex: 0,
		name: 'Numerical value',
        type: 'value',
		position: 'left'
    },{
		gridIndex: 0,
		name: 'Residual',
        type: 'value',
		position: 'right'
    },{
		gridIndex: 1,
		name: 'Standardized residual',
        type: 'value',
		position: 'left'
    },{
		gridIndex: 1,
		name: 'Relative deviation',
        type: 'value',
		//offset:40,
		position: 'right'
    }],
	series: [
        {
            name:'Actual value',
            type:'line',
			xAxisIndex: 0,
			yAxisIndex: 0,
            data:[0,1,2,3]
        },
        {
            name:'Model value',
            type:'line',
			xAxisIndex: 0,
			yAxisIndex: 0,
            data:[0,1,2,3]
        },
        {
            name:'Residual',
            type:'line',
			xAxisIndex: 0,
			yAxisIndex: 1,
            data:[0,1,2,3]
        },
        {
            name:'Standardized residual',
            type:'line',
			xAxisIndex: 1,
			yAxisIndex: 2,
            data:[0,1,2,3]
        },
        {
            name:'Relative deviation',
            type:'line',
			xAxisIndex: 1,
			yAxisIndex: 3,
            data:[0,1,2,3]
        }
    ]
};
function refreshEcharts_his()
{	
	var steptype='';
	if(TAG.TagType=='BOOL'){
		steptype='end';
	}
	hisChart.setOption(his_option, true);
	hisChart.setOption({
		title: {
			text: TAG.Name+'-Historical data and increment'
		},
		xAxis: [
		{
			data: HIS_TIME
		}],
		series: [
		{
			step: steptype,
			data: HIS_DATA
		},{
			step: steptype,
			data: HIS_INCREMENT_DATA
		}]
	});
	$("#Echarts_His").show();
};

function refreshEcharts_hisInterval()
{	
	var legends=[];
	for( var yl=0;yl < TAG_Y.length;yl++){
        if(TAG_Y.hasOwnProperty(yl))
		    legends[yl]=TAG_Y[yl].Name;
    }
    for( var xl=0;xl < HIS_INTERVAL_DATA.length;xl++){
        if(TAGS_SERIAL.hasOwnProperty(xl))
            legends[xl+yl]=TAGS_SERIAL[xl].Name;
    }
	hisIntervalChart.setOption(hisInterval_option, true);
    hisIntervalChart.setOption({
		title: {
            text: 'Data comparison chart'
        },
		legend: {
			data: legends
		},
		xAxis: [
		{
            data: HIS_INTERVAL_TIME
        }],
        series:function(){
	    	var serie=[];
		    for( var i=0;i < (xl+yl);i++){
				var steptype='';
				var tagname;
				var datas=[];
				if(i<yl){
                    if(TAG_Y.hasOwnProperty(i)){//存在
                        if(TAG_Y[i].TagType=='BOOL'){
                            steptype='end';
                        }
                        tagname=TAG_Y[i].Name;
                        datas=HIS_INTERVAL_DATA_Y;
                    }
				}else{
					var j=i-yl;
					if(TAGS_SERIAL.hasOwnProperty(j)){//存在
						if(TAGS_SERIAL[j].TagType=='BOOL'){
							steptype='end';
						}
						tagname=TAGS_SERIAL[j].Name;
						datas=HIS_INTERVAL_DATA[j];
					}
				}
			   	var item={
					name:tagname,
					type:'line',
					markPoint : {
						data : [
							{type : 'max', name : 'Max'},
							{type : 'min', name : 'Min'}
						]
					},
					markLine : {
						data : [
							{type : 'average', name : 'Average'}
						]
					},
					step: steptype,
					data: datas
				};
				serie.push(item);
			};
	    	//console.log(serie);
	    	return serie;
	    }()
    });
	$("#Echarts_HisIntervalSerial").show();
	$("#HisSerialRemark").show();
};

function refreshRegressionEcharts()
{
    trendChart.setOption({
		xAxis: [
		{
            data: HIS_INTERVAL_TIME
        },{
            data: HIS_INTERVAL_TIME
        }],
        series: [
		{
            data: REGRESSION_RES.Ys
        },{
            data: REGRESSION_RES.YEst
        },{
            data: REGRESSION_RES.Residual
        },{
            data: REGRESSION_RES.StdRes
        },{
            data: REGRESSION_RES.RelDev
        }]
    });

	scatterChart.setOption({
        series: [
		{
            data: REG_Y_LIMIT
        },{
            data: REG_Y_UP_SIGMA
        },{
            data: REG_Y_BELOW_SIGMA
        },{
            data: REGRESSION_RES.Yscatter
        }]
    });

	$("#Echarts_Trend").show();
	$("#Echarts_Scatter").show();
}

if (hisInterval_option && typeof hisInterval_option === "object") {
    hisIntervalChart.setOption(hisInterval_option, true);
	document.getElementById("Echarts_HisIntervalSerial").style.display="none";//隐藏
	document.getElementById("HisSerialRemark").style.display="none";
}

if (his_option && typeof his_option === "object") {
    hisChart.setOption(his_option, true);
	document.getElementById("Echarts_His").style.display="none";//隐藏
} 

if (scatter_option && typeof scatter_option === "object") {
    scatterChart.setOption(scatter_option, true);
	document.getElementById("Echarts_Scatter").style.display="none";//隐藏
}
if (trend_option && typeof trend_option === "object") {
    trendChart.setOption(trend_option, true);
	document.getElementById("Echarts_Trend").style.display="none";//隐藏
}