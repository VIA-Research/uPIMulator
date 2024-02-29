package reg_descriptor

type SrcRegDescriptor struct {
	gp_reg_descriptor *GpRegDescriptor
	sp_reg_descriptor *SpRegDescriptor
}

func (this *SrcRegDescriptor) InitGpRegDescriptor(gp_reg_descriptor *GpRegDescriptor) {
	this.gp_reg_descriptor = gp_reg_descriptor
	this.sp_reg_descriptor = nil
}

func (this *SrcRegDescriptor) InitSpRegDescriptor(sp_reg_descriptor *SpRegDescriptor) {
	this.gp_reg_descriptor = nil
	this.sp_reg_descriptor = sp_reg_descriptor
}

func (this *SrcRegDescriptor) IsGpRegDescriptor() bool {
	return this.gp_reg_descriptor != nil
}

func (this *SrcRegDescriptor) IsSpRegDescriptor() bool {
	return this.sp_reg_descriptor != nil
}

func (this *SrcRegDescriptor) GpRegDescriptor() *GpRegDescriptor {
	return this.gp_reg_descriptor
}

func (this *SrcRegDescriptor) SpRegDescriptor() *SpRegDescriptor {
	return this.sp_reg_descriptor
}
