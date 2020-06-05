package readline

func (reader *Reader) AppendHistory(prompt string) {
	reader.History = append(reader.History, prompt)
}
