/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2023/2/20 17:17:48                           */
/*==============================================================*/


drop table if exists User;

drop table if exists comment;

drop table if exists favorite;

drop table if exists video;

/*==============================================================*/
/* Table: User                                                  */
/*==============================================================*/
create table User
(
   username             varchar(32),
   password             varchar(32),
   userid               int not null,
   follow_count         int,
   follower_count       int,
   avatar               varchar(60),
   background_image     varchar(60),
   signature            varchar(60),
   total_favorited      varchar(60),
   work_count           int,
   favorite_count       int,
   primary key (userid)
);

/*==============================================================*/
/* Table: comment                                               */
/*==============================================================*/
create table comment
(
   commentid            int not null,
   userid               int not null,
   videoid              int not null,
   content              varchar(50),
   create_data          varchar(20),
   primary key (commentid)
);

/*==============================================================*/
/* Table: favorite                                              */
/*==============================================================*/
create table favorite
(
   likeid               int not null,
   userid               int not null,
   videoid              int not null,
   primary key (likeid)
);

/*==============================================================*/
/* Table: video                                                 */
/*==============================================================*/
create table video
(
   videoid              int not null,
   userid               int not null,
   play_url             varchar(60),
   cover_url            varchar(60),
   comment_count        int,
   title                varchar(30),
   vedio_favorite_count int,
   primary key (videoid)
);

alter table comment add constraint FK_user_comment foreign key (userid)
      references User (userid) on delete restrict on update restrict;

alter table comment add constraint FK_vedio_comment foreign key (videoid)
      references video (videoid) on delete restrict on update restrict;

alter table favorite add constraint FK_user_like foreign key (userid)
      references User (userid) on delete restrict on update restrict;

alter table favorite add constraint FK_vedio_like foreign key (videoid)
      references video (videoid) on delete restrict on update restrict;

alter table video add constraint FK_publish foreign key (userid)
      references User (userid) on delete restrict on update restrict;

