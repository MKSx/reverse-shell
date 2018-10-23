package handler

type ProcessOutput struct {
	Result  []byte
	Process *Process
}

type ProcessTerminated struct {
	Reason  []byte
	Process *Process
}
