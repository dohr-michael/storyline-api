package relation

type Type string

const (
	CreatedBy = Type("createdBy")
	IsOwnerOf = Type("isOwnerOf")
	BelongsTo = Type("belongsTo")
)
