package grpc

// type healthRes struct {
// 	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
// }

type runRes struct {
	Computation string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}