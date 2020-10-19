# MicEngine

* 更新日志

|版本      |时间      |更新内容      |
|:--------|:-------- |:--------   |
|V1.1.2039b|2020-10-14|1.主程序中增加对看门狗的监视,发现看门狗停止后自动启动之|
|V1.1.2038|2020-10-12|1.修正报表层级树中存在的权限不正确的问题;<br>2.修正多国语言页面中的部分翻译问题|
|V1.1.2037|2020-09-21|1.增加中英文切换功能; <br>2.增加看门狗程序micdog.exe; <br>3.修正读取长时间庚顿数据时存在的bug; <br>**升级注意**:<br>confg文件夹下新增micdog.ini文件|
|V1.1.2036|2020-09-09|1.增加自定义的首页面。<br> 2.主菜单增加英文项。<br>**升级注意:**<br>1)config.ini配置文件中需要增加FirstPage配置项。<br>2)sys_meneu数据库表需要增加name_eng字段|
|V1.1.2035|2020-08-28|1.增加对系统时序数据表中数据的统计接口和对应的巡检数据的统计结果窗口;<br>2.增加本地Sqlite运行心跳记录|
|V1.1.2034|2020-08-17|增加写庚顿数据库快照和历史数据的API接口|
|V1.1.2033|2020-08-14|1.手动同步taglist到庚顿数据库时,同期保存庚顿数据库的标签点ID到taglist的alias_tag_id字段。<br>2.自适应平台和Grafana的内外网路径。<br>**升级注意:**<br>平台登录计算服务时api接口中需要增加tip参数,用于传送登录浏览器的host。<br>需要在sys_dictionary中添加 dic_catalog_id分别为13、26、27,dictionary_code分别为1、2的记录，如图![](http://47.93.34.228/showdoc/server/../Public/Uploads/2020-08-14/5f363d7041131.png)|
|V1.1.2032|2020-08-03|1.优化对SQL脚本,解决SQL查询结果只有一行数据时返回为0的bug。<br>2.修正读取庚顿历史数据时如果没有读取到数据不能正确报错的bug。<br>3.修正读取庚顿等间隔历史数据时获取的数据间隔与设定间隔有细微差别的bug|
|V1.1.2031|2020-07-17|1.增加KPI数据可视化页面。<br>2.优化对仿SQL脚本的识别。<br>3.程序内部结构优化。<br>4.快照页面增加报警信息显示。|
|V1.1.2030|2020-07-09|1.优化统计指标，在历史数据中增加四分位统计.<br>2.修正因为BOOL变量设定了最大最小值而引起的不统计开停时间的问题。 <br>3.优化报表在线预览中的对Excel公式的支持。<br>4.需要在配置文件cfig.ini中增加报表Excel公式计算深度选项[ExcelFormulaCalcDeep = 5]|
|V1.1.2029|2020-07-03|优化报表在线预览的时对Excel公式的支持,可以支持除了IF()条件语句之外的大部分常用Excel计算公式 |
|V1.1.2028|2020-06-29|1.为庚顿数据库增加连接池,每次并发最大50连接。<br>2.对实时数据统计时,只取质量码为GOOD的数据。<br>3.增加操作庚顿数据库的API接口。<br>4.需要在配置文件cfig.ini中增加庚顿数据库连接池大小选项[GoldenCennectPool = 50]|
|V1.1.2027|2020-06-19|1. 在KPI计算和报表计算中增加自定义SQL脚本功能，可以全平台自由查询数据。 <br>2. 对实时数据统计时,只取最大值和最小值之间的数据。如果没有设置最大值和最小值，则不进行此项筛选|
|V1.1.2026|2020-06-16|1.解决了无法往庚顿中写入非ASCII码字符的问题。<br>2.读取庚顿数据库的历史数据时，如果时间范围超过8个小时，则分为多个线程并发读取，以提高读取效率。<br>3.需要在配置文件cfig.ini中增加读取历史数据时的分段时间长度(小时)选项[MaxHisSliceTime = 8] |
|V1.1.2025|2020-06-10|1.增加模拟量数据报警功能模块。<br>2.调整化验数据类型，新增product_lab,并将lab标签与product_lab标签等效处理。<br>**升级注意:**1.升级后首次运行需要将配置文件中的[createtable]设置为true;<br>2.需要在配置文件cfig.ini中增加报警使能选项[AlarmEnable = true / false] <br>3.需要将平台数据库ore_process_d_taglist表中的 category、location1~location5字段的类型由varchar更改为int|
|V1.1.2024|2020-06-08|1.配合Windows计划任务，可以进行自我检测，实现故障时自启动；<br>2.修正首页面Seesion问题;<br>3.修正date函数中的参数bug,并增加可选项<br>4.允许跨域请求|
|V1.1.2023|2020-05-29|报表在线预览模式增加对Excel公式的有限度支持|
|V1.1.2022|2020-05-28|1.报表增加调试模式,在调试模式下错误的公式和错误原因会列在报表中,费调试模式下错误公式和错误信息不会显示在报表中,错误信息保存在日志中。<br>2.项目信息中增加报表授权数量显示。<br>3.**升级注意:**升级后首次运行需要将配置文件中的[createtable]设置为true|
|V1.1.2021|2020-05-25|1.优化报表在线预览界面，仿Excel界面。<br>2.增加报表模板在线编辑功能,仿Excel编辑方式。<br>3.修正srtd函数中某一段时间没有数据而停止KPI计算的bug|
|V1.1.2020|2020-05-21|1.优化树形菜单，增加更新、展开、折叠、隐藏、搜索等功能<br>2.增加快照中错误信息提醒功能，点击[#Error]可以看具体错误信息<br>3.修正Grafana数据源中描述信息中存在的Bug|
|V1.1.2019|2020-05-19|将庚顿驱动程序由WebAPI更新为C接口<br>**升级注意:**<br>1.需要在calc_kpi_engine_config中将庚顿WebAPI地址更换为庚顿数据库的直接地址和端口,并设置用户名和密码。<br>2.不支持庚顿变量名中的非英文字符，庚顿变量名中有非英文字母的，需要更改|
|V1.1.2018|2020-05-12|增加报表运算功能<br>升级注意:1.需要在配置文件cfig.ini中增加报表使能选项[ReportEnable = true];<br>2.升级后首次运行需要将配置文件中的[createtable]设置为true;<br>3.需要在URL管理中为报表相关菜单授权才可看到报表相关页面。|
|V1.1.2017|2020-04-29|1.进入页面时校验该用户是否有该页面的授权。<br>2.修正快照页面加载显示，切换层级时显示数据加载提示。<br>3.更换为蓝色LOGO|
|V1.1.2016|2020-04-21|增加日志查询分析模块 |
|V1.1.2015|2020-04-19|1.增加巡检模块<br>2.增加项目信息页面|
|V1.1.2014|2020-04-15|增加质检和物耗模块|
|V1.1.2013|2020-03-27|1.大版本更新,增加数据分析可视化页面|
|V1.1.2013beta|2020-03-16|1.大版本更新,增加数据分析可视化页面，但尚未测试完成;<br>修正srtd中diff计算中存在的错误|
|V1.0.2012|2020-03-02|修正载入process2指标的时候存在的bug|
|V1.0.2011|2020-03-01|1.增加每条KPI自定义base_time和shift_hour参数功能;<br>2.时间参数中增加班、日、月、季度、年起始时间描述|
|V1.0.2010|2020-2-24 |脚本中增加t_add(basetime,duration)时间偏移函数|
|V1.0.2009|2020-2-18 |增加在process_taglist中定义的变量，但数据存储在平台sys_real_data数据表中的过程变量计算（process2）|
|V1.0.2008|2020-2-17 |1.在计算服务引擎配置表中增加记录软件版本的功能;<br>2.增加校验脚本的API接口;<br>3.修正srtd中求差时指定时间点没有数据的时候不会自动后移时间段的问题|
|V1.0.2007|2020-2-14 |优化峰谷值计算,增加taglist.user_int2作为是否删除负数的参数,增加taglist.user_real2作为是"增量临界值"参数|
|V1.0.2006|2020-2-13 |1.增加默认平台数据库地址和端口地址;<br>2.引擎配置表中增加心跳计数器|
|V1.0.2005|2020-2-9  |增加sample化验KPI指标计算|
|V1.0.2004|2020-2-5  |增加峰谷值计算|
|V1.0.2003|2020-1-15 |1.修正串行计算没有数据时停滞不前的问题； <br>2.调整offset_minute的用法|
|V1.0.2002|2020-1-13 |增加平台动态数据表计算功能|
|V1.0.2001|2020-1-12 |初版发布|