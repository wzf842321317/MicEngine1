//<!--Echarts脚本文件-->
var pageEcharts = echarts.init(document.getElementById("DataAnalysMsg"));
//------------饼图----------------------
echarts_option = {
    color: ['#ff7f50','#87cefa','#da70d6','#32cd32','#6495ed',
    '#ff69b4','#ba55d3','#cd5c5c','#ffa500','#40e0d0',
    '#1e90ff','#ff6347','#7b68ee','#00fa9a','#ffd700',
    '#6699FF','#ff6666','#3cb371','#b8860b','#30e0e0'],
};
function echartsDataInit(){
    var i=1;
    var datas=new Array();
    var item={
        name: 'Source of visit',
        type: 'pie',
        radius: '40%',
        center: ['25%', '25%'],
        data: [],
        emphasis: {
            itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
        }
    };
    for(let k in SOURCE_ANALYSE){
        var d={value:SOURCE_ANALYSE[k],name:k};
        item.data.push(d);
    }
    datas.push(item);

    for(let k in SOURCE_ANALYSE){
        i++;
        var point=['25%', '25%'];
        switch(i){
            case 2:
                point[0]='25%';
                point[1]='75%';
                break;
            case 3:
                point[0]='75%';
                point[1]='25%';
                break;
            case 4:
                point[0]='75%';
                point[1]='75%';
                break;
            default:
                break;
        }
        var it={
            name: k,
            type: 'pie',
            radius: '40%',
            center: point,
            data: [],
            emphasis: {
                itemStyle: {
                    shadowBlur: 10,
                    shadowOffsetX: 0,
                    shadowColor: 'rgba(0, 0, 0, 0.5)'
                }
            }
        };

        switch(k){
        case "Management platform":
            for(let k2 in PLAT_ANALYSE){
                var d={value:PLAT_ANALYSE[k2],name:k2};
                it.data.push(d);
            }
            break;
        case "Mobile app":
            for(let k3 in APP_ANALYSE){
                var d={value:APP_ANALYSE[k3],name:k3};
                it.data.push(d);
            }
            break;
        case "Analysis system":
            for(let k4 in ENG_ANALYSE){
                var d={value:ENG_ANALYSE[k4],name:k4};
                it.data.push(d);
            }
            break;
        }
        datas.push(it);
    }
    return datas;
}
function refreshPageEcharts(sites,datas){
    pageEcharts.setOption(echarts_option, true);
    pageEcharts.setOption({
        title: {
            text: 'Access information statistics',
            subtext: 'Access source and URL',
            left: 'center'
        },
        tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b} : {c} ({d}%)'
        },
        series: echartsDataInit()
    });    
    $("#DataAnalysMsg").show();
}
if (echarts_option && typeof echarts_option === "object") {
    pageEcharts.setOption(echarts_option, true);
	document.getElementById("DataAnalysMsg").style.display="none";//隐藏
};