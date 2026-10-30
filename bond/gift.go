package bond

// Talk is a chat: once a day, it makes them a little fonder.
const Talk = 10

// Give is a gift on day of something they like as much as liked, which
// may be below 0 for what they dislike: one gift a day.
func (b *Bond) Give(liked, day int) error {
	if b.gifted == day+1 {
		return ErrToday
	}
	b.gifted = day + 1
	b.fonder(liked)
	return nil
}

// Chat is talking with them on day: fonder by Talk, once a day.
func (b *Bond) Chat(day int) error {
	if b.talked == day+1 {
		return ErrToday
	}
	b.talked = day + 1
	b.fonder(Talk)
	return nil
}
