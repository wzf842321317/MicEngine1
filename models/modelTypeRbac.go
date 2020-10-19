package models

//用户表
type SysUser struct {
	Id               int64               //
	Name             string              //姓名
	CardNo           string              //身份证号
	Username         string              //用户名
	Password         string              //密码
	Telephone        string              //手机号
	HeadPic          string              //头像
	PositionId       int64               //岗位ID
	DepartmentId     int64               //部门ID
	Email            string              //邮箱
	Sex              int64               //性别：0-女，1-男
	Status           int64               //状态：1-有效，0-无效
	UserStatus       int64               //人员状态：1-在职，2-离职，3-挂职，4-借调
	Remark           string              //备注
	Operator         int64               //操作者id
	OperateTime      string              `orm:"type(datetime)"` //操作时间
	OperateIp        string              //操作者的IP地址
	Createtime       string              `orm:"type(datetime)"` //记录的创建时间
	UserFrom         int64               //用户来自哪儿
	UserType         int64               //用户类型
	Qq               string              //qq 号
	Weixin           string              //微信号
	Permission       int64               //权限值
	ValidTime        string              `orm:"type(datetime)"` //有效时间
	LoginType        int64               //登录类型
	EmailCode        string              //通过邮箱找回密码时的验证码
	NewEmail         string              //新邮箱号
	EmailStatus      int64               //邮箱状态
	RealValidStatus  int64               //实名认证状态
	Roles            []*SysRole          `orm:"rel(m2m);rel_through(github.com/bkzy-wangjp/MicEngine/models.SysRoleUser)"` //设置与角色的多对多关系
	GoodsDataCreates []*GoodsConsumeInfo `orm:"reverse(many)"`
	GoodsDataUpdates []*GoodsConsumeInfo `orm:"reverse(many)"`
	CheckPlans       []*CheckPlan        `orm:"reverse(many)"` //设置一对多反向关系
	SysLogs          []*SysLog           `orm:"reverse(many)"` //设置一对多反向关系
	SampleToLabs     []*SamplelistToLab  `orm:"reverse(many)"` //设置一对多反向关系
	LabResults       []*LabAnaResultTsd  `orm:"reverse(many)"` //设置一对多反向关系
}

//角色表
type SysRole struct {
	Id          int64            //
	RoleName    string           //角色名称
	Remarkes    string           //角色说明
	Users       []*SysUser       `orm:"reverse(many)"`                                                                   //设置与用户表的反向多对多关系
	Permissions []*SysPermission `orm:"rel(m2m);rel_through(github.com/bkzy-wangjp/MicEngine/models.SysRolePermission)"` //设置与权限表的反向多对多关系
}

//权限表
type SysPermission struct {
	Id              int64                  //
	PermissionName  string                 //权限名称
	Remarkes        string                 //备注
	PermissionValue int64                  //权限值
	PermissionType  int64                  //1菜单权限，2.操作权限 ,3.数据权限,4.实时监控URL,5.本地可视化,6.APP首页应用权限,7报表配置信息管理
	Roles           []*SysRole             `orm:"reverse(many)"`                                                                           //设置与角色的多对多关系
	Menus           []*SysMenu             `orm:"rel(m2m);rel_through(github.com/bkzy-wangjp/MicEngine/models.SysMenuPermission)"`         //设置与菜单的多对多关系
	Middles         []*MsMiddleCorrelation `orm:"rel(m2m);rel_through(github.com/bkzy-wangjp/MicEngine/models.SysMiddlePermission)"`       //设置与菜单的多对多关系
	Monitors        []*RealTimeMonitor     `orm:"rel(m2m);rel_through(github.com/bkzy-wangjp/MicEngine/models.RealTimeMonitorPermission)"` //设置与监控菜单的多对多关系
	Reports         []*CalcKpiReportList   `orm:"rel(m2m);rel_through(github.com/bkzy-wangjp/MicEngine/models.SysReportPermission)"`       //设置与菜单的多对多关系
}

//菜单表
type SysMenu struct {
	Id              int64            //ID
	Name            string           //菜单名称
	NameEng         string           //英文菜单名
	Pid             int64            //父级菜单
	Level           int64            //层级
	Seq             int64            //排序号
	Remark          string           //备注
	Url             string           //actionURL
	HasChild        int64            //是否包含下级，1-包含，0-不包含
	CreateUser      string           //创建人
	CreateTime      string           `orm:"type(datetime)"` //创建时间
	PermissionValue int64            //所需权限值
	State           int64            //1有效 0无效
	MenuType        int64            //平台菜单为1,本地可视化菜单为2
	IconUrl         string           //菜单图标url
	Permissions     []*SysPermission `orm:"reverse(many)"` //设置与权限表的反向多对多关系
}

//用户角色关系表
type SysRoleUser struct {
	Id   int64    //ID
	User *SysUser `orm:"rel(fk)"`
	Role *SysRole `orm:"rel(fk)"`
}

//角色权限关系表
type SysRolePermission struct {
	Id         int64          //ID
	Permission *SysPermission `orm:"rel(fk)"`
	Role       *SysRole       `orm:"rel(fk)"`
}

//权限菜单关系表
type SysMenuPermission struct {
	Id         int64          //ID
	Permission *SysPermission `orm:"rel(fk)"`
	Menu       *SysMenu       `orm:"rel(fk)"`
}

//中间表权限关系表
type SysMiddlePermission struct {
	Id         int64                //ID
	Permission *SysPermission       `orm:"rel(fk)"`
	Middle     *MsMiddleCorrelation `orm:"rel(fk)"`
}

//权限-监控关系表
type RealTimeMonitorPermission struct {
	Id         int64            //ID
	Permission *SysPermission   `orm:"rel(fk)"`
	Monitor    *RealTimeMonitor `orm:"rel(fk);column(table_id)"`
}

//权限-报表关系表
type SysReportPermission struct {
	Id         int64              //ID
	Permission *SysPermission     `orm:"rel(fk)"`
	Report     *CalcKpiReportList `orm:"rel(fk)"`
}
