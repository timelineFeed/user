-- 用户信息表
create table `user `(
                        `id` bigint unsigned auto_increment,
                        `name` varchar(30)  not null comment '用户昵称',
                        `password` varchar(128) not null comment '用户密码hash',
                        `telephone` varchar(11) not null  default '' comment '用户电话号码',
                        `email` varchar(40) not null default '' comment '用户邮箱号',
                        `image` text not null default '' comment '用户头像',
                        `status` tinyint unsigned not null default 0 comment '状态，10-删除',
                        `created_at` datetime not null  comment '创建时间',
                        `updated_at` datetime not null comment '更新时间',
                        primary key (`id`)
)engine=InnoDB default charset=utf8;
