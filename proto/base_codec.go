package proto


type baseCodec struct {
	order      *ExtendedByteOrder
	capability *Capability
}

func newBaseCodec() *baseCodec {
	c := Capability(0)
	return &baseCodec{
		order: &ExtendedByteOrder{LittleEndian},
		capability: &c}
}

func (this *baseCodec)Capability() Capability {
	return *this.capability
}
func (this *baseCodec)SetCapability(c Capability) {
	*this.capability = c
}
func (this *baseCodec) HasCapability(c Capability) bool {
	return this.capability.Has(c)
}
