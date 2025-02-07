package liberdatabase

import "gorm.io/gorm"

// LiberTextDocumentCharacter represents a character in a text document
type LiberTextDocumentCharacter struct {
	gorm.Model
	TextDocumentId int64  `gorm:"column:document_id;not null"`
	Character      string `gorm:"column:character;not null"`
	Count          int64  `gorm:"column:count;not null"`
}

func (LiberTextDocumentCharacter) TableName() string {
	return "public.liber_text_document_characters"
}
