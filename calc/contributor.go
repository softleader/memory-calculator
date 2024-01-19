package calc

type Contributor interface {

	// Contribute 在執行底層 libjbm 前, 所要貢獻的內容, 這會在執行前先被呼叫, 通常是一些參數的設定等
	Contribute() error
}
