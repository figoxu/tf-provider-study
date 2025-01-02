# tf-provider-study

# gaussdb 的 schema声明

`tf-provider-study/terraform-provider-huaweicloud/huaweicloud/services/gaussdb/resource_huaweicloud_gaussdb_opengauss_instance.go:331`

## 资源维护的相关接口

terraform-provider-huaweicloud/huaweicloud/services/gaussdb/resource_huaweicloud_gaussdb_opengauss_instance.go:369

# GaussDB(for openGauss) 操作系统类型
## 默认操作系统
GaussDB(for openGauss) 默认使用 CentOS 作为基础操作系统。
## 主要特点
- 基于 CentOS 7.X 版本定制
- 预装所有数据库运行依赖
- 系统镜像由华为云维护
- 用户无法修改操作系统类型
## 使用限制
- 创建实例时无需指定 os_type 参数
- 不支持自定义操作系统类型
- 不支持直接访问底层操作系统
- 所有运维操作需通过华为云管理接口执行

