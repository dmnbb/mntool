package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"mntool/idmaker"
	"strings"
)

func main() {
	idmakerDemo()
}

//create database mndb
func idmakerDemo() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mndb")
	if err != nil {
		fmt.Printf("[ERRO] sql.Open failed err=%v\n", err)
		return
	}

	_, err = db.Exec("CREATE TABLE `id_maker` (\n  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',\n  `biz_id` tinyint(4) NOT NULL DEFAULT '0' COMMENT '业务id',\n  `next_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '下一个id',\n  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',\n  `update_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '更新时间',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `uniq_biz_id` (`biz_id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='id生成器';")
	if err != nil && !strings.Contains(err.Error(), "Error 1050:") {
		fmt.Printf("[ERRO] create table id_maker failed err=%v\n", err)
		return
	}

	_, err = db.Exec("INSERT INTO id_maker\n(biz_id, next_id, create_time, update_time)\nVALUES(1, 1000, unix_timestamp(), unix_timestamp());")
	if err != nil && !strings.Contains(err.Error(), "Error 1062:") {
		fmt.Printf("[ERRO] insert id_maker failed err=%v\n", err)
		return
	}

	if err = idmaker.New(1, db, 10); err != nil {
		fmt.Printf("[ERRO] idmaker.New failed err=%v\n", err)
		return
	}

	for i := 0; i < 3; i++ {
		id, err := idmaker.GetId(1)
		if err != nil {
			fmt.Printf("[ERRO] idmaker.GetId failed err=%v\n", err)
			return
		}
		fmt.Printf("[INFO] GetId=%d\n", id)
	}
}
