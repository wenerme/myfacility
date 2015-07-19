/*
MySQL Protocol speaks in GO

Packet Manipulation
Use Read and Writer to do basic protocol stuff, use Proto to work with packet.

	Get(&value,&value,&value...)
	Get(&value,ProtoType)
	Get(&value,StrVar,n)
	Get(&value,Int,n)
	Get(n, IgnoreByte)
	Get(&value,reflect.Kind)

Change Get to Put for Writer.When Put, the pointer is unnecessary but acceptable.
When use ProtoType to access integer must use unsigned type.Most of time, use Reader and Writer is just a method name change.
This is how to decode and encode OK packet, it's easy to read and write, all packets follow this pattern.

	func (p *OKPack) Read(c Proto) {
		c.Get(&p.Header, &p.AffectedRows, IntEnc, &p.LastInsertId, IntEnc)
		if c.HasCap(CLIENT_PROTOCOL_41) {
			c.Get(&p.Warnings, &p.Status)
		} else if c.HasCap(CLIENT_TRANSACTIONS) {
			c.Get(&p.Status)
		}

		if c.HasCap(CLIENT_SESSION_TRACK) {
			c.Get(&p.Info)
			if p.Status.Has(SERVER_SESSION_STATE_CHANGED) {
				c.Get(&p.SessionState)
			}
		} else {
			c.Get(&p.Info, StrEof)
		}
	}
	func (p *OKPack) Write(c Proto) {
		c.Put(&p.Header, &p.AffectedRows, IntEnc, &p.LastInsertId, IntEnc)
		if c.HasCap(CLIENT_PROTOCOL_41) {
			c.Put(&p.Warnings, &p.Status)
		} else if c.HasCap(CLIENT_TRANSACTIONS) {
			c.Put(&p.Status)
		}

		if c.HasCap(CLIENT_SESSION_TRACK) {
			c.Put(&p.Info)
			if p.Status.Has(SERVER_SESSION_STATE_CHANGED) {
				c.Put(&p.SessionState)
			}
		} else {
			c.Put(&p.Info, StrEof)
		}
	}

When use Get/Put without explicit type, will use following type map:

	*uint -> IntEnc
	*int -> IntEnc will convert to/from uint
	*[]byte -> StrVar
	*string -> StrVar will convert to/from []byte

Other type will just use binary.Read/Write(r/w, binary.LittleEndian, v).
For string,[]byte type, Reader and Writer behave is not same,Reader will use binary.Read,Writer will use StrVar.
*/
package proto
