## 变更日志

### 2026-05-03 Round 2: 数据完整性与功能完善

**novel 模块:**
- dao/directory_dao.go: 新增 GetChaptersByNovelID 函数
- service/directory_service.go: UpdateChapter 时调用 countWordsFromHTML 重新计算 WordCount
- service/sync_helper.go: 新增 countWordsFromHTML() 中英文混合字数统计函数
- service/recent_service.go: LogRecentAccess 无章节时回退到记录小说级别访问
- handler/sync_handler.go: 新增 GET /api/novels/:novelId/updates 同步端点
- handler/export_handler.go: 新增 POST /api/novels/:novelId/export 导出端点
- handler/node_crud_handler.go: 新增 settings/custom-plot/custom-analysis/custom-others 单项 CRUD
- router/novel_router.go: 注册新增路由

**ai 模块:**
- dao/conversation_dao.go: 新增 FindByID/Update/Delete 方法
- service/conversation_service.go: 新增 UpdateTitle/SaveMessages/Delete
- handler/chat_handler.go: 新增 UpdateConversation/DeleteConversation handler
- router/ai_router.go: 注册 PUT/DELETE /api/ai/chat/conversations/:id
- model/common.go: StreamResponse 添加 Event 字段
- provider/openai.go, claude.go, gemini.go: 所有 chunk 添加 Event 字段

**settings 模块:**
- handler/stub_handlers.go: system/themes/settings, usage-logs, privacy 的 stub 端点
- router/settings_router.go: 注册 system/usage-logs/privacy 路由

### 2026-05-03 Phase 3: 修复后端致命 Bug

**修改文件:**

1. **dao/recent_dao.go** — 修复 LogOrUpdateRecentActivity 的 TOCTOU 竞态条件
   - 修改前: 先 `First` 查询再 `Create` 插入，并发时可能产生重复记录
   - 修改后: 使用 GORM `Clauses(clause.OnConflict{...})` 原子化 upsert，一个 SQL 完成创建或更新

2. **service/history_service.go** — 添加 RestoreVersion 的 documentID 校验
   - 修改前: 任何属于用户的版本可以回滚到任意文档上，造成跨文档数据污染
   - 修改后: 校验 `versionToRestore.DocumentID == documentID`，防止跨文档恢复

3. **dao/novel_dao.go** — 两处修复
   - PermanentlyDeleteNovelByID: 添加 `Preload("DerivedContents").Preload("Notes")` 确保 GORM 能级联删除派生内容和笔记
   - UpdateNovelJSONField: 添加字段白名单，防止通过 `fieldName` 参数注入更新非预期的数据库列

4. **ai/model/common.go** — StreamResponse 添加 Event 字段
   - 修改前: 响应格式 `{content, done}`，前端 `aiApi.ts` 检查 `parsedData.event === 'chunk'` 永远不成立
   - 修改后: 添加 `Event string` 字段，前端可正确区分 chunk/done/error 事件

5. **ai/provider/openai.go, claude.go, gemini.go** — 所有 StreamResponse 实例添加 Event 字段
   - 修改前: 所有 content chunk 缺少 Event 字段
   - 修改后: content chunk 设置 `Event: "chunk"`，done 设置 `Event: "done"`，error 设置 `Event: "error"`

6. **ai/handler/chat_handler.go** — 错误响应添加 Event 字段

---

1:这里完成小说数据的持久
2：小说数据的增删改查
