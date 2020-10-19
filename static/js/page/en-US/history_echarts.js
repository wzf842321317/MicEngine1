var his_dom = document.getElementById("Echarts_His");
var hisChart = echarts.init(his_dom); 

var hisInterval_dom = document.getElementById("Echarts_HisInterval");
var hisIntervalChart = echarts.init(hisInterval_dom); 

var hisGroup_dom = document.getElementById("Echarts_HisGroup");
var hisGroupChart = echarts.init(hisGroup_dom);
//------------单变量历史趋势图(通用型)----------------------
 his_option = {
    title: {
        text: 'Original historical data'
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
        boundaryGap: false,
        data:  [1,2,3,4,5,6,7,8]
    }],
    yAxis: [{
		gridIndex: 0,
		name: 'Data',
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

//------------描述统计分组----------------------
 hisGroup_option = {
    title: {
        text: 'Statistical grouping data'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross'
        }
    },

	color:['red','green','blue','#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83',  '#ca8622', '#bda29a','#6e7074', '#546570', '#c4ccd3'],
    legend: {
        data: ['Data distribution map']
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
		name: 'Section',
        boundaryGap: false,
        data:  [1,2,3,4,5,6,7,8]
    }],
    yAxis: [{
		gridIndex: 0,
		name: 'Proportion (%)',
        type: 'value',
		position: 'left'
    }],
	series: [
        {
            name:'Data distribution map',
            type:'bar',
            data: [1,2,3,4,5,6,7,8]
        }
    ]
}; 

function refreshEcharts_his()
{	
	var steptype='';
	if(TAG.TagType=='BOOL'){
		steptype='end';
	}
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
	var steptype='';
	if(TAG.TagType=='BOOL'){
		steptype='end';
	}
    hisIntervalChart.setOption({
		title: {
            text: TAG.Name+'-Aggregate data at equal intervals'
        },
		xAxis: [
		{
            data: HIS_INTERVAL_TIME
        }],
        series: [
		{
			step: steptype,
            data: HIS_INTERVAL_DATA
        }]
    });
	$("#Echarts_HisInterval").show();
	$("#Echarts_HisGroup").show();
};

function refreshEcharts_hisGroup()
{
    //console.log(CHART_Q1,CHART_Mean,CHART_Median,CHART_Q3);
    hisGroupChart.setOption({
		title: {
            text: TAG.Name+'-Data distribution map'
        },
		xAxis: [
		{
            data: HIS_SUM_GROUP_KEY
        }],
        series: [
		{
            data: HIS_SUM_GROUP_VAL,
            markArea: {
                data: [[{
                    name: 'Quartile range',
                    xAxis: CHART_Q1
                }, {
                    xAxis: CHART_Q3
                }]]
            },
            markLine: {
                lineStyle: {
                    type: 'solid',
                    color: '#c23531'
                },
                data: [
                    { xAxis: CHART_Median,name:'Median' }
                ]
            }
        }]
    });
};

if (hisInterval_option && typeof hisInterval_option === "object") {
    hisIntervalChart.setOption(hisInterval_option, true);
	document.getElementById("Echarts_HisInterval").style.display="none";//隐藏
}
if (hisGroup_option && typeof hisGroup_option === "object") {
    hisGroupChart.setOption(hisGroup_option, true);
	document.getElementById("Echarts_HisGroup").style.display="none";//隐藏
}
if (his_option && typeof his_option === "object") {
    hisChart.setOption(his_option, true);
	document.getElementById("Echarts_His").style.display="none";//隐藏
} 
