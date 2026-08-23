# 转发会话取消复现

一次投递中出现被拒绝的日志后，转发会话会把取消状态留给下一批投递。

先用同一会话提交 ID 为 `broken` 的拒绝记录，再提交 ID 为 `next` 的正常记录。

第二批返回 `context canceled`，没有正常回执；目标检查报出 `following delivery inherited the cancelled batch`。
