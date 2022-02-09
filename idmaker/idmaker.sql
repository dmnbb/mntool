
-- 建表
CREATE TABLE `id_maker` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `biz_id` tinyint(4) NOT NULL DEFAULT '0' COMMENT '业务id',
  `next_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下一个id',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_biz_id` (`biz_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='id生成器';

-- 业务注册
INSERT INTO id_maker
(biz_id, next_id, create_time, update_time)
VALUES(1, 1000, unix_timestamp(), unix_timestamp());
