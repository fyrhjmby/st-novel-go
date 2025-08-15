package main

import (
	"sync"
	"time"
)

// AIDataStore 封装了AI模块所需的所有数据
type AIDataStore struct {
	ChatConversations      map[string]Conversation
	chatConversationsMutex sync.RWMutex
}

// initAIData 初始化AI模块的模拟数据
func initAIData() {
	chatConversations := map[string]Conversation{
		"conv-1": {
			ID: "conv-1", Title: "角色性格讨论", Summary: "关于卡尔文的内心世界...", CreatedAt: time.Now().Format(time.RFC3339),
			Messages: []ChatMessage{
				{ID: "msg-1", Role: "user", Content: "我想讨论一下卡尔文这个角色的内心世界。在面对跃迁点时，他应该有什么样的心理活动？", Timestamp: time.Now().Add(-2 * time.Hour).Format(time.RFC3339)},
				{ID: "msg-2", Role: "ai", Content: `<p class="mb-3">当然，这是一个很好的切入点。卡尔文在面对跃迁点时的心理活动可以从多个层面展现：</p><ol class="list-decimal list-inside space-y-2"><li><strong>希望与恐惧的交织</strong>：跃迁点代表着回家的可能，但也意味着巨大的未知风险。可以描写他手心出汗，但眼神却异常明亮。</li><li><strong>孤独感的顶峰</strong>：在做出这个重大决定时，他比任何时候都更加意识到自己的孤独。没有战友，只有AI。</li><li><strong>对过去的眷恋</strong>："回家"这个词会触发他的记忆闸门，可以闪回一些与家人或恋人的片段。</li></ol>`, Timestamp: time.Now().Add(-2*time.Hour + 1*time.Second).Format(time.RFC3339)},
			},
		},
		"conv-2": {ID: "conv-2", Title: "情节灵感", Summary: "跃迁点的科学原理...", CreatedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339), Messages: []ChatMessage{}},
	}

	AIStore = AIDataStore{
		ChatConversations:      chatConversations,
		chatConversationsMutex: sync.RWMutex{},
	}
}
