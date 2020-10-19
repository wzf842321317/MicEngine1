//<!--Echarts脚本文件-->
var Echarts_dom = document.getElementById("PageEcharts");
var DataEcharts = echarts.init(Echarts_dom);

//------------单变量历史趋势图(时间对齐)----------------------
	echarts_option = {
    title: {
        text: '等间隔聚合历史数据'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross'
        }
    },

	color:['blue','red','green','#c23531','#2f4554', '#61a0a8', '#d48265', '#91c7ae','#749f83',  '#ca8622', '#bda29a','#6e7074', '#546570', '#c4ccd3'],
    legend: {
        data: ['等间隔聚合数据']
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
		name: '时间',
        boundaryGap: false,
        data:  [1,2,3,4,5,6,7,8]
    }],
    yAxis: [{
		gridIndex: 0,
		name: '数值',
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
            name:'等间隔聚合数据',
            type:'line',
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
			step: '',
            data: [1,2,3,4,5,6,7,8]
        }
    ]
}; 

function refreshEcharts_hisInterval()
{	
	var legends=[];
	for( var i=0;i < HIS_INTERVAL_DATA.length;i++){
		legends[i]=TAGS_SERIAL[i].Name;
	}
    DataEcharts.setOption({
		title: {
            text: '数据对比图'
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
		    for( var i=0;i < HIS_INTERVAL_DATA.length;i++){
				var steptype='';
				if(TAGS_SERIAL[i].TagType=='BOOL'){
					steptype='end';
				}
			   	var item={
					name:TAGS_SERIAL[i].Name,
					type:'line',
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
					step: steptype,
					data: HIS_INTERVAL_DATA[i]
				};
				serie.push(item);
			};
	    	console.log(serie);
	    	return serie;
	    }()
    });
	$("#PageEcharts").show();
};

if (echarts_option && typeof echarts_option === "object") {
    DataEcharts.setOption(echarts_option, true);
	document.getElementById("PageEcharts").style.display="none";//隐藏
};