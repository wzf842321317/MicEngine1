//<!--Echarts脚本文件-->
var Echarts_dom = document.getElementById("PageEcharts");
var DataEcharts = echarts.init(Echarts_dom);
echarts_option={
    title: {
        text: 'KPI计算结果数据'
        //subtext: '纯属虚构'
    },
    tooltip: {
        trigger: 'axis'
    }
};
function refreshEcharts_Bar()
{	
    DataEcharts.setOption({
            title: {
                text: 'KPI计算结果数据'
                //subtext: '纯属虚构'
            },
            tooltip: {
                trigger: 'axis',
                formatter: '{a} <br/>{b} : {c}'
            },
            legend: {
                left:'center',
                data: ['KPI']
            },
            toolbox: {
                show: true,
                x:'left',
                feature: {
                    dataView: {show: true, readOnly: false},
                    magicType: {show: true, type: ['line', 'bar']},
                    restore: {show: true},
                    saveAsImage: {show: true}
                }
            },
            calculable: true,
            xAxis: [
                {
                    type: 'category',
                    show:true,
                    data: DATA_TIME
                }
            ],
            yAxis: [
                {
                    type: 'value',
                    show:true
                }
            ],
            series: [
                {
                    name: 'KPI',
                    type: 'bar',
                    data: DATA_VALUE,
                    markPoint: {
                        data: [
                            {type: 'max', name: '最大值'},
                            {type: 'min', name: '最小值'}
                        ]
                    },
                    markLine: {
                        data: [
                            {type: 'average', name: '平均值'}
                        ]
                    }
                }
            ]    
    });
	$("#PageEcharts").show();
};

function refreshEcharts_Pie(titlestr,)
{	
    DataEcharts.setOption({
        title: {
            text: titlestr,
            //subtext: '纯属虚构',
            left: 'center'
        },
        tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b} : {c} ({d}%)'
        },
        toolbox: {
            show: false
        },
        legend: {
            //orient: 'vertical',
            bottom: 10,
            left: 'center',
            data: PIE_LEGEND
        },
        xAxis: [
            {
                show:false
            }
        ],
        yAxis: [
            {
                show:false
            }
        ],
        series: [
            {
                name: titlestr,
                type: 'pie',
                radius: [0,'60%'],
                center: ['50%', '50%'],
                data: PIE_DATA,
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
	$("#PageEcharts").show();
};

if (echarts_option && typeof echarts_option === "object") {
    DataEcharts.setOption(echarts_option, true);
	document.getElementById("PageEcharts").style.display="none";//隐藏
};