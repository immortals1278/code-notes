# 设计
在 FetchPending 时，增加 gracePeriod 参数，给热路径 gracePeriod 的时间将已发送的信息从 pending 改成 published
# 语法
```golang
//使用 PostgreSQL 的 COPY FROM 协议高性能批量插入数据库
_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"orders"},// 表名
		[]string{"id", "user_id", "symbol", "side", "price", "quantity", "filled_quantity", "status"},// 确定字段
		pgx.CopyFromRows(rows),// 插入的数据，在rows中，rows是二维切片
	)
```

```golang
//rows是db读到的数据
//for rows.Next() 循环里,每调一次 Next() 把游标移到下一行,然后 Scan 读取当前行的数据。
//Scan 等价于："把这 10 列的值，按顺序塞进 m 的 10 个字段里"。
for rows.Next() {
		m := &Message{}
		err := rows.Scan(
			&m.ID, &m.AggregateID, &m.AggregateType,
			&m.Topic, &m.PartitionKey, &m.Payload,
			&m.Status, &m.RetryCount, &m.CreatedAt, &m.PublishedAt,
		)
}
```

`WHERE id = ANY($1)` 等价于 WHERE id IN (?, ?, ?)

```golang
// WithTx 在单个 db 事务中执行回调函数，供 Worker 持有 SKIP LOCKED 鎖直到批次處理完成。
//SKIP LOCKED:如果某行已被其他事务锁住,直接跳过而不是等待
func (r *Repository) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		txCtx := context.WithValue(ctx, db.TxKey, tx)
		return fn(txCtx)
	})
}
```

`time.Now().UnixMilli()` time.TIME的方法，以毫秒返回现在的时间

```golang
rows, err := r.getExecutor(ctx).Query(ctx, query, x)
defer rows.Close()// 关闭数据库结果集（rows），防止数据泄露
```

`rows.Err()` 用于检查遍历 rows 过程中是否发生错误