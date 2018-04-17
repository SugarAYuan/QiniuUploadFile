## 这是一个七牛云的上传视频文件的服务

- config.json为配置文件 如果不需要数据库请注释相关入库程序后重新编译
- videoTest 为编译好的可执行文件可直接运行
- 命令行提供连个参数 
    1. 一个是要上传的视频文件url 程序将会自动抓取视频资源放到本地
    2. 第二个参数是视频保存的文件名
    
> 数据库

```mysql
CREATE TABLE `video` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '订单id',
  `url` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '地址',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `delete_time` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci
```