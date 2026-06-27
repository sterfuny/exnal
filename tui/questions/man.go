package questions

import (
	tea "charm.land/bubbletea/v2"
)
// Question 所有题目类型都要实现的接口
type Question interface {
    // Render 渲染题目（包括当前输入/选择状态）
    Render() string
    
    // HandleKey 处理键盘输入，返回 (是否完成, 是否退出)
    HandleKey(msg tea.KeyMsg) (done bool, quit bool)
    
    // GetQuestionText 获取题目文字
    GetQuestionText() string
    
    // GetAnswer 获取答案（用于最终汇总）
    GetAnswer() any
    
    // IsDone 是否已答完
    IsDone() bool
}
