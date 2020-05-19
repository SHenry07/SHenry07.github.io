```sql
select count(distinct 学号) as 学生人数 
from score;
```

```mysql
select sum(成绩)
from score
where 课程号 = '0002';
```

```mysql
select count(教师号)
from teacher
where 教师姓名 like '孟%';
```

```mysql
/*
分析思路
select 查询结果 [课程ID：是课程号的别名,最高分：max(成绩) ,最低分：min(成绩)]
from 从哪张表中查找数据 [成绩表score]
where 查询条件 [没有]
group by 分组 [各科成绩：也就是每门课程的成绩，需要按课程号分组];
*/
select 课程号,max(成绩) as 最高分,min(成绩) as 最低分
from score
group by 课程号;
```