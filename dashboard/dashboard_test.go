package dashboard

import (
	"fmt"
	"regexp"
	"testing"
)

func TestAdd(t *testing.T) {
	a := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Awesome go-echarts</title>
    <script src="https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js"></script>
    <link href="https://go-echarts.github.io/go-echarts-assets/assets/bulma.min.css" rel="stylesheet">
</head>

<body>
<div class="select" style="margin-right:10px; margin-top:10px; position:fixed; right:10px;"></div>

<div class="container">
    <div class="item" id="EqqeQmdiTaJS"
         style="width:100%;height:330px;"></div>
</div>
<script type="text/javascript">
    "use strict";
    let myChart_EqqeQmdiTaJS = echarts.init(document.getElementById('EqqeQmdiTaJS'), "white");
    let option_EqqeQmdiTaJS = {
        title: {},
        tooltip: {},
        legend: {"left":"1%","top":"2%",},
        toolbox: {"show":true,"feature":{"saveAsImage":{},"dataZoom":{},"dataView":{},"restore":{}}},
        xAxis: [{"data":["2019-12-20\n21:04:48","2019-12-20\n21:09:48","2019-12-20\n21:14:48","2019-12-20\n21:19:48","2019-12-20\n21:24:48","2019-12-20\n21:29:48"],"splitArea":{"show":false,},"splitLine":{"show":false,}}],
        yAxis: [{"axisLabel":{"show":true,"formatter":"{value} %"},"min":0,"max":100,"splitArea":{"show":false,},"splitLine":{"show":true,"lineStyle":{"type":"dashed"}}}],
        series: [
        {"name":"system","type":"line","data":[5.25,1.82,1.63,4.92,3.88,2.8],"label":{"show":true},"emphasis":{"label":{"show":false},},"markLine":{"label":{"show":false}},"markPoint":{"label":{"show":false}},},
        {"name":"user","type":"line","data":[23.33,4.22,3.38,22.73,11.49,9.16],"label":{"show":true},"emphasis":{"label":{"show":false},},"markLine":{"label":{"show":false}},"markPoint":{"label":{"show":false}},},
        {"name":"iowait","type":"line","data":[],"label":{"show":true},"emphasis":{"label":{"show":false},},"markLine":{"label":{"show":false}},"markPoint":{"label":{"show":false}},},
        ],
        color: ["#c23531","#2f4554","#61a0a8","#d48265","#91c7ae","#749f83","#ca8622","#bda29a","#6e7074","#546570","#c4ccd3"],
    };
    myChart_EqqeQmdiTaJS.setOption(option_EqqeQmdiTaJS);
</script>

<style>
    .container {margin-top:30px; display: flex;justify-content: center;align-items: center;}
    .item {margin: auto;}
</style>
</body>
</html>`

	reg, _ := regexp.Compile(`{([\s\S]*)};`)
	fmt.Println(reg.FindString(a))
}
