package statements

import (
	"crf-mold/base"
	"time"
)

type CreateStatementVO struct {
	Order     int       `json:"order"`                       //
	CreatedAt time.Time `json:"createdAt"`                   //
	UpdatedAt time.Time `json:"updatedAt"`                   //
	Id        string    `json:"id"`                          // 主键id
	CourseId  string    `json:"courseId" binding:"required"` // 课程id
	SoundMark string    `json:"soundMark"`                   // 音标
	Chinese   string    `json:"chinese"`                     // 中文
	English   string    `json:"english" binding:"required"`  // 英文
	Pid       string    `json:"pid"`                         // 父id
} // @name CreateStatementVO

type UpdateStatementVO struct {
	Order     int       `json:"order"`                 //
	CreatedAt time.Time `json:"createdAt"`             //
	UpdatedAt time.Time `json:"updatedAt"`             //
	Id        string    `json:"id" binding:"required"` // 主键id
	Video     string    `json:"video"`                 //
	CourseId  string    `json:"courseId"`              // 课程id
	SoundMark string    `json:"soundMark"`             // 音标
	Chinese   string    `json:"chinese"`               // 中文
	English   string    `json:"english"`               // 英文
	Pid       string    `json:"pid"`                   // 父id
} // @name UpdateStatementVO

type ListStatementVO struct {
	CourseId string `json:"courseId" binding:"required"` // 课程id
} // @name ListStatementVO

type PageStatementVO struct {
	base.PageVO
	Order     int       `json:"order"`                        //
	CreatedAt time.Time `json:"createdAt"`                    //
	UpdatedAt time.Time `json:"updatedAt"`                    //
	Id        string    `json:"id" binding:"required"`        // 主键id
	Video     string    `json:"video"`                        //
	CourseId  string    `json:"courseId" binding:"required"`  // 课程id
	SoundMark string    `json:"soundMark" binding:"required"` // 音标
	Chinese   string    `json:"chinese" binding:"required"`   // 中文
	English   string    `json:"english" binding:"required"`   // 英文
	PId       string    `json:"pid" binding:"required"`       // 父id
} // @name PageStatementVO

type PageStatementOutputVO struct {
	Order     int       `json:"order"`                        //
	CreatedAt time.Time `json:"createdAt"`                    //
	UpdatedAt time.Time `json:"updatedAt"`                    //
	Id        string    `json:"id" binding:"required"`        // 主键id
	Video     string    `json:"video"`                        //
	CourseId  string    `json:"courseId" binding:"required"`  // 课程id
	SoundMark string    `json:"soundMark" binding:"required"` // 音标
	Chinese   string    `json:"chinese" binding:"required"`   // 中文
	English   string    `json:"english" binding:"required"`   // 英文
	PId       string    `json:"pid" binding:"required"`       // 父id
} // @name PageStatementOutputVO

type StatementOutVO struct {
	Order     int       `json:"order"`                        //
	CreatedAt time.Time `json:"createdAt"`                    //
	UpdatedAt time.Time `json:"updatedAt"`                    //
	Id        string    `json:"id" binding:"required"`        // 主键id
	Video     string    `json:"video"`                        //
	CourseId  string    `json:"courseId" binding:"required"`  // 课程id
	SoundMark string    `json:"soundMark" binding:"required"` // 音标
	Chinese   string    `json:"chinese" binding:"required"`   // 中文
	English   string    `json:"english" binding:"required"`   // 英文
	PId       string    `json:"pid" binding:"required"`       // 父id
} // @name StatementOutVO

type SplitStatementVO struct {
	CourseId string `json:"courseId" binding:"required"` // 课程id
	//Statement string `json:"statement" binding:"required"` // 句子
	StatementIds []string `json:"statementIds"` // 句子ids
} // @name StatementOutVO

type DeleteStatementVO struct {
	CourseId string `json:"courseId" binding:"required"` // 课程id
	//Statement string `json:"statement" binding:"required"` // 句子
	StatementIds []string `json:"statementIds"` // 句子ids
} // @name DeleteStatementVO

type DeepSeekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekChatCompletionsResponse struct {
	Choices []struct {
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
		LogProbs     string  `json:"logprobs"`
		Message      Message `json:"message"`
	} `json:"choices"`
	Created           int    `json:"created"`
	ID                string `json:"id"`
	Model             string `json:"model"`
	Object            string `json:"object"`
	SystemFingerprint string `json:"system_fingerprint"`
	Usage             struct {
		CompletionTokens      int `json:"completion_tokens"`
		PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
		PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
		PromptTokens          int `json:"prompt_tokens"`
		TotalTokens           int `json:"total_tokens"`
	} `json:"usage"`
}

type Statement struct {
	Chinese   string `json:"chinese"`
	English   string `json:"english"`
	Order     int    `json:"order"`
	Soundmark string `json:"soundmark"`
}

type SentenceData struct {
	Sentence   string      `json:"sentence"`
	Statements []Statement `json:"statements"`
}
