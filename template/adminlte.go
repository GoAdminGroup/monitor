package template

import "html/template"

type Adminlte struct{}

func (*Adminlte) GetDashboardStyle() string {
	return `<script>
	$("body").addClass("sidebar-collapse")
	$(".zoom-in-btn").on('click', function (event) {
		$(this).parent().parent().parent().addClass("zoom-in-container")
		let container = $(this).parent().next().children().children()
		container.attr("data-raw-height", container.height())
		container.css("height", "680px");
		let chartID = "myChart_" + container.attr("id")
		eval(chartID + ".resize()")
		$(this).hide()
		$(this).next().show()
	});
	$(".zoom-out-btn").on('click', function (event) {
		$(this).parent().parent().parent().removeClass("zoom-in-container")
		let container = $(this).parent().next().children().children() 
		container.css("height", container.attr("data-raw-height"));
		let chartID = "myChart_" + container.attr("id")
		eval(chartID + ".resize()")
		$(this).hide()
		$(this).prev().show()
	});
	$(".refresh-btn").on('click', function (event) {
		let container = $(this).parent().next().children().children() 
		let chartID = "myChart_" + container.attr("id")
		$.ajax({
			url: "%s/refresh/" + $(this).attr("data-chart-id") + "/chart/dashboard/%s",
			data: {
				query_time_start: Date.parse(new Date(new Date().getTime() - 1 * 60 * 60 * 1000))/1000,
				query_time_end: Date.parse(new Date())/1000
			},
			type: "POST",
			success: function(data) {
				if (data.code == 0) {
					eval(chartID + ".setOption(" + data.data + ")")
				}
			}
		});
	});
</script>
<style>
	.echarts-container{margin-top:0px;}
	.zoom-in-container{position:absolute;width:93.5%%;height:82%%;z-index:999;}
	.row {margin-right: 0;margin-left: 0;}
	.col-lg-1, .col-lg-10, .col-lg-11, .col-lg-12, .col-lg-2, .col-lg-3, .col-lg-4, .col-lg-5, .col-lg-6, .col-lg-7, .col-lg-8, .col-lg-9, .col-md-1, .col-md-10, .col-md-11, .col-md-12, .col-md-2, .col-md-3, .col-md-4, .col-md-5, .col-md-6, .col-md-7, .col-md-8, .col-md-9, .col-sm-1, .col-sm-10, .col-sm-11, .col-sm-12, .col-sm-2, .col-sm-3, .col-sm-4, .col-sm-5, .col-sm-6, .col-sm-7, .col-sm-8, .col-sm-9, .col-xs-1, .col-xs-10, .col-xs-11, .col-xs-12, .col-xs-2, .col-xs-3, .col-xs-4, .col-xs-5, .col-xs-6, .col-xs-7, .col-xs-8, .col-xs-9{padding-right: 6px;padding-left: 6px;}
	.box {margin-bottom: 6px;border: 1px solid #d2d6de;}
	.box:hover {border: 1px solid #5298b3;}
</style>`
}

func (*Adminlte) GetToolBar(interval string) template.HTML {
	if interval == "" || (interval != "15" && interval != "30" && interval != "60" && interval != "120") {
		interval = "60"
	}

	intervalTimes := map[string]string{
		"15":  "15000",
		"30":  "30000",
		"60":  "60000",
		"120": "120000",
	}

	indexs := map[string]int{
		"15":  3,
		"30":  2,
		"60":  1,
		"120": 0,
	}

	choosedClass := []string{"", "", "", ""}
	choosedClass[indexs[interval]] = " choosed-btn"

	chooseBtn := template.HTML(`
<style>
.choose-time-btn {
	padding: 4px 8px 4px 8px;
	background-color: #ffffff;
	color: #7d7c7c;
	margin-right: 10px;
	font-size: 13px;
	cursor: pointer;
	width: 50px;
	text-align: center;
	float: right;
}
.choose-time-btn:hover {
	background-color: #d4d2d2;
}
.choose-time-btn.choosed-btn {
	background-color: #546478;
	color: #fff;
}
.breadcrumb {
	display: none;
}
</style>
<div style="position: absolute;right: 18px;top: 67px;">
	<div class="choose-time-btn` + choosedClass[0] + `" data="120">2分钟</div>
	<div class="choose-time-btn` + choosedClass[1] + `" data="60">1分钟</div>
	<div class="choose-time-btn` + choosedClass[2] + `" data="30">30秒</div>
	<div class="choose-time-btn` + choosedClass[3] + `" data="15">15秒</div>
</div>
<script>
$(".choose-time-btn").on("click", function (e) {
	if (!$(this).hasClass("choosed-btn")) {
		$(".choose-time-btn.choosed-btn").removeClass("choosed-btn");
		$(this).addClass("choosed-btn");
		location.href = "?interval=" + $(this).attr("data");
	}
});
window.setTimeout(function(){
	$.pjax.reload('#pjax-container');
}, ` + intervalTimes[interval] + `);
</script>
`)
	return chooseBtn
}

func (*Adminlte) GetGraphBtn() string {
	return `<div class="zoom-in-btn" style="cursor: pointer;color: #a7a7a7;
float: right;margin-right: 1%%;"><i class="fa fa-arrows-alt"></i></div><div class="zoom-out-btn" style="display:none;cursor: pointer;color: #a7a7a7;
float: right;margin-right: 1%%;"><i class="fa fa-compress"></i></div><div class="refresh-btn" data-chart-id="%d" style="cursor: pointer;color: #a7a7a7;
float: right;margin-right: 2.2%%;"><i class="fa fa-refresh"></i></div>`
}

func (*Adminlte) GetSingleStatBtn() string {
	return `<div class="zoom-in-btn" style="cursor: pointer;color: #a7a7a7;
float: right;margin-right: 1%%;"><i class="fa fa-arrows-alt"></i></div><div class="zoom-out-btn" style="display:none;cursor: pointer;color: #a7a7a7;
float: right;margin-right: 1%%;"><i class="fa fa-compress"></i></div><div class="refresh-btn" data-chart-id="%d" style="cursor: pointer;color: #a7a7a7;
float: right;margin-right: 2.2%%;"><i class="fa fa-refresh"></i></div>`
}

func (*Adminlte) GetSingleStatContent() string {
	return `<div style="height:%dpx;line-height: %dpx;font-size: 2.5em;font-weight:bold;text-align:center;color: %s;">%s</div>`
}
